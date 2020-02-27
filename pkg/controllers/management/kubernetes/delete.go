/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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

		// @TODO this needs to be reverted once the fix in placed UI
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   object.Spec.Provider.Group,
			Kind:    "GKE",
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
		if err := kubernetes.DeleteIfExists(ctx, a.mgr.GetClient(), &v1.Secret{
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
