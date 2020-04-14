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

package eks

import (
	"context"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the aws eks cluster
func (t *eksCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to delete eks cluster")

	// @step: first we need to check if we have access to the credentials
	resource := &eks.EKS{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	result, err := func() (reconcile.Result, error) {
		operations := []controllers.EnsureFunc{
			t.EnsureDeletionStatus(resource),
			t.EnsureNodeGroupsDeleted(resource),
			t.EnsureDeletion(resource),
			t.EnsureSecretDeletion(resource),
		}

		for _, handler := range operations {
			result, err := handler(ctx)
			if err != nil {
				return reconcile.Result{}, err
			}
			if result.Requeue || result.RequeueAfter > 0 {
				return result, nil
			}
		}

		// @cool we can remove the finalizer now
		finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
		if err := finalizer.Remove(resource); err != nil {
			logger.WithError(err).Error("removing the finalizer from eks resource")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the eks cluster")
		resource.Status.Status = corev1.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the resource status")

		return reconcile.Result{}, err
	}

	return result, nil
}
