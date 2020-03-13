/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package eks

import (
	"context"
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the gke cluster
func (t *eksCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Info("attempting to delete eks cluster")

	// @step: first we need to check if we have access to the credentials

	resource := &eksv1alpha1.EKS{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	requeue, err := func() (bool, error) {
		creds, err := t.GetCredentials(ctx, resource, request.NamespacedName.Name)
		if err != nil {
			return false, err
		}

		// @step: create a cloud client for us
		client, err := NewClient(creds, resource)
		if err != nil {
			return false, err
		}

		// @step: check if the cluster exists and if so we wait or the operation or the exit
		found, err := client.Exists()
		if err != nil {
			return false, fmt.Errorf("checking if cluster exists: %s", err)
		}

		// @step: lets update the status of the resource to deleting
		if resource.Status.Status != corev1.DeleteStatus {
			resource.Status.Status = corev1.DeleteStatus

			return true, nil
		}

		if found {
			_, err := client.Delete()
			// TODO know which errors re should retry / reque
			return false, err
		}

		return false, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the cluster")

		return reconcile.Result{}, err
	}
	if requeue {
		if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource status")

			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	if err := finalizer.Remove(resource); err != nil {
		logger.WithError(err).Error("removing the finalizer")

		return reconcile.Result{}, err
	}

	logger.Info("successfully deleted the cluster")

	return reconcile.Result{}, nil
}
