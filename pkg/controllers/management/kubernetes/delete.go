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

package kubernetes

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ComponentKubernetesCleanup = "Kubernetes Clean-up"
)

// Delete is responsible for deleting any bindings which were created
func (a k8sCtrl) Delete(ctx context.Context, object *clustersv1.Kubernetes) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      object.Name,
		"namespace": object.Namespace,
	})
	logger.Debug("attempting to delete the cluster from the api")

	original := object.DeepCopy()

	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				a.EnsureDeleteStatus(object),
				a.EnsureServiceDeletion(object),
				a.EnsureSecretDeletion(object),
			},
		)
	}()
	if err != nil {
		logger.WithError(err).Error("trying to delete the kubernetes resource")
		object.Status.Status = corev1.DeleteFailedStatus
	}

	// @step: update the status of the resource
	if err := a.mgr.GetClient().Status().Patch(ctx, object, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the status of the resource")

		return reconcile.Result{}, err
	}

	// @cool we can remove the finalizer now
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)
	if err := finalizer.Remove(object); err != nil {
		logger.WithError(err).Error("removing the finalizer from eks resource")

		return reconcile.Result{}, err
	}

	return result, nil
}
