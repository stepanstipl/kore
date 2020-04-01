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
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ComponentClusterDelete is the name of the cluster deletion component
	ComponentClusterDelete = "Cluster Deletor"
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

		object.Status.Status = corev1.DeleteStatus

		// @check if a we are backed by a cloud provider
		if kore.IsProviderBacked(object) {
			logger := logger.WithFields(log.Fields{
				"group":     object.Spec.Provider.Group,
				"kind":      u.GetKind(),
				"name":      u.GetName(),
				"namespace": u.GetNamespace(),
				"version":   object.Spec.Provider.Version,
			})

			// @step: retrieve the cloud provider resource from the api
			ref := types.NamespacedName{Namespace: u.GetNamespace(), Name: u.GetName()}

			if err := a.mgr.GetClient().Get(ctx, ref, u); err != nil {
				if !kerrors.IsNotFound(err) {
					logger.WithError(err).Error("trying to retrieve the cluster resource from the api")

					return reconcile.Result{}, err
				}
			} else {
				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterDelete,
					Message: "attempting to clean up in-cluster kubernetes resources",
				})

				// Assumption: we only clean up everything if we own the provider

				// First delete all namespaces to ensure this will work
				// @step: retrieve the provider credentials secret
				token, err := controllers.GetConfigSecret(context.Background(),
					a.mgr.GetClient(),
					object.Namespace,
					object.Name)
				if err != nil {
					if !kerrors.IsNotFound(err) {
						logger.WithError(err).Error("unable to obtain credentials secret")

						object.Status.Components.SetCondition(corev1.Component{
							Name:    ComponentClusterDelete,
							Message: "Unable obtain cluster access cluster credentials",
							Detail:  err.Error(),
							Status:  corev1.FailureStatus,
						})

						return reconcile.Result{}, err
					}
				}

				if token != nil {
					// @step: create a client for the remote cluster
					client, err := kubernetes.NewRuntimeClientFromConfigSecret(token)
					if err != nil {
						logger.WithError(err).Error("trying to create client from credentials secret")

						object.Status.Components.SetCondition(corev1.Component{
							Name:    ComponentClusterDelete,
							Message: "Unable to access cluster using provided cluster credentials",
							Detail:  err.Error(),
							Status:  corev1.FailureStatus,
						})

						return reconcile.Result{}, err
					}
					err = CleanupKoreCluster(ctx, client)
					if err != nil {
						logger.WithError(err).Error("trying to clean up cluster resources")

						object.Status.Components.SetCondition(corev1.Component{
							Name:    ComponentClusterDelete,
							Message: "Unable to delete all cluster namespaces",
							Detail:  err.Error(),
							Status:  corev1.FailureStatus,
						})

						return reconcile.Result{}, err
					}
				}

				object.Status.Components.SetCondition(corev1.Component{
					Name:    ComponentClusterDelete,
					Message: "Waiting for cloud provider to be deleted",
				})

				// @step: should we delete the cloud provider if there is one
				if a.Config().EnableClusterDeletion && u.GetDeletionTimestamp() == nil {
					logger.Info("attempting to delete the cloud provider from the api")

					// @step: we should attempt to delete the cloud provider
					if err := a.mgr.GetClient().Delete(ctx, u); err != nil {
						logger.WithError(err).Error("trying delete the cloud cluster")

						object.Status.Components.SetCondition(corev1.Component{
							Name:    ComponentClusterDelete,
							Message: "Failed trying to delete the cloud provider",
							Detail:  err.Error(),
						})

						return reconcile.Result{}, err
					}

					return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
				}

				if a.Config().EnableClusterDeletionBlock {
					// @check if cloud provider is still being deleted we wait
					return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
				}
			}
			logger.Debug("attempting to delete the kubernetes credential")
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
