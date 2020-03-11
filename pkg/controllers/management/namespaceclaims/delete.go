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
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	original := resource.DeepCopy()

	// @step: check if we are the current finalizer
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	result, err := func() (reconcile.Result, error) {
		logger.Debug("deleting the namespaceclaim from the cluster")

		// @step: create a client from the cluster secret
		client, err := controllers.CreateClientFromSecret(context.Background(), a.mgr.GetClient(),
			resource.Spec.Cluster.Namespace, resource.Spec.Cluster.Name)
		if err != nil {
			logger.WithError(err).Error("trying to create kubernetes client from secret")

			return reconcile.Result{}, err
		}

		// @step: update the status of the resource
		if resource.Status.Status != corev1.DeleteStatus {
			resource.Status.Status = corev1.DeleteStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: delete the namespace
		if err := kubernetes.DeleteIfExists(context.Background(), client, &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: resource.Name},
		}); err != nil {
			logger.WithError(err).Error("trying to delete the namespace contained in the claim")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Failed to delete namespaceclaim",
		}}

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		// @step: update the status of the resource
		if err := a.mgr.GetClient().Status().Patch(context.Background(), resource, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource status")

			return reconcile.Result{}, err
		}

		return result, nil
	}

	// @step: remove the finalizer if one and allow the resource it be deleted
	if err := finalizer.Remove(resource); err != nil {
		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Failed to remove the finalizer",
		}}

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
