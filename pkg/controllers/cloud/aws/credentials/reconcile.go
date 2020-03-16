/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package credentials

import (
	"context"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t awsCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to reconcile aws credentials")

	resource := &eks.EKSCredentials{}
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
			logger.WithError(err).Error("trying to create aws permissions client")

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
			logger.Warn("awscredentials not verified")
		}

		return reconcile.Result{}, err
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the aws credentials")

		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Either the credentials are invalid or we've encountered an error verifying",
		}}
	}

	// @step: update the status of the resource
	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the aws credentials resource status")

		return reconcile.Result{}, err
	}

	return result, err
}
