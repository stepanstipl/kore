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
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	result, err := func() (reconcile.Result, error) {
		// @step: we need to grab the cloud provider and check if it's deletion
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   object.Spec.Provider.Group,
			Kind:    object.Spec.Provider.Kind,
			Version: object.Spec.Provider.Version,
		})
		u.SetName(object.Name)
		u.SetNamespace(object.Namespace)

		object.Status.Status = corev1.DeletingStatus

		// @check if a we are backed by a cloud provider
		if kore.IsProviderBacked(object) {
			// Assumption: we only clean up everything if we own the provider
			if object.Spec.Provider.Kind == "EKS" {
				if err := a.EnsureResourceDeletion(context.Background(), object); err != nil {
					return reconcile.Result{}, err
				}
			}

		}

		// @step: we should delete the secert from api
		if err := kubernetes.DeleteIfExists(ctx, a.mgr.GetClient(), &configv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      object.Name,
				Namespace: object.Namespace,
			},
		}); err != nil {
			log.WithError(err).Error("trying to delete the secret from api")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to delete the kubernetes resource")
	}
	if err == nil {
		if result.RequeueAfter <= 0 && !result.Requeue {
			if err := finalizer.Remove(object); err != nil {
				log.WithError(err).Error("trying to remove the finalizer from resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	// @step: update the status of the resource
	if err := a.mgr.GetClient().Status().Patch(ctx, object, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the status of the resource")

		return reconcile.Result{}, err
	}

	return result, nil
}
