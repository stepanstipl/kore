/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package namespaceclaims

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for removig the namespace claim any remote configuration
func (a *nsCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.Name,
		"namespace": request.Namespace,
	})
	logger.Debug("attempting to delete the resource")

	// @step: retrieve the resource from the api
	resource := &clustersv1.NamespaceClaim{}
	if err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, resource); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// @step: check if we are the current finalizer
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	err := func() error {
		logger.Debug("deleting the namespaceclaim from the cluster")

		// @step: create a client from the cluster secret
		client, err := controllers.CreateClientFromSecret(context.Background(), a.mgr.GetClient(),
			resource.Spec.Cluster.Namespace, resource.Spec.Cluster.Name)
		if err != nil {
			logger.WithError(err).Error("trying to create kubernetes client from secret")

			return err
		}

		// @step: delete the namespace
		if err := kubernetes.DeleteIfExists(context.Background(), client, &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: resource.Name},
		}); err != nil {
			logger.WithError(err).Error("trying to delete the namespace contained in the claim")

			return err
		}

		return nil
	}()
	if err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed to delete namespaceclaim",
		}}

		return reconcile.Result{}, err
	}

	// @step: remove the finalizer if one and allow the resource it be deleted
	if err := finalizer.Remove(resource); err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed to remove the finalizer",
		}}

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
