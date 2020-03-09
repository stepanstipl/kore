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

package gkecredentials

import (
	"context"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t gkeCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to reconcile gke credentials")

	resource := &gke.GKECredentials{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	result, err := func() (reconcile.Result, error) {
		var verified bool
		resource.Status.Conditions = []corev1.Condition{}
		resource.Status.Verified = &verified

		// @step: set the status to pending if none set
		if resource.Status.Status == "" {
			resource.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: create the client to verify the permissions
		client, err := NewClient(resource)
		if err != nil {
			logger.WithError(err).Error("trying to create gcp permissions client")

			return reconcile.Result{}, err
		}

		// @step: verify the credentials
		verified, err = client.HasRequiredPermissions()
		if err != nil {
			return reconcile.Result{}, err
		}
		resource.Status.Verified = &verified
		resource.Status.Status = corev1.SuccessStatus

		if resource.Status.Verified == nil || !*resource.Status.Verified {
			logger.Warn("gke credentials not verified")
		}

		return reconcile.Result{}, err
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gke credentials")

		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Either the credentials are invalid or we've encountered an error verifying",
		}}
	}

	// @step: update the status of the resource
	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the gke credentials resource status")

		return reconcile.Result{}, err
	}

	return result, err
}
