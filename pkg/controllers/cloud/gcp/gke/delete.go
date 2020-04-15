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

package gke

import (
	"context"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the gke cluster
func (t *gkeCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Info("attempting to delete gke cluster")

	// @step: first we need to check if we have access to the credentials
	resource := &gke.GKE{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	result, err := func() (reconcile.Result, error) {
		creds, err := t.GetCredentials(ctx, resource, request.Namespace)
		if err != nil {
			return reconcile.Result{}, err
		}

		// @step: create a cloud client for us
		client, err := NewClient(creds, resource)
		if err != nil {
			return reconcile.Result{}, err
		}

		// @step: we need to retrieve the current state
		cluster, found, err := client.GetCluster(ctx)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the state of the cluster")

			return reconcile.Result{}, err
		}

		if found {
			switch cluster.Status {
			case "PROVISIONING", "RECONCILING":
				return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
			case "ERROR", "RUNNING":
				// @step: lets update the status of the resource to deleting
				if resource.Status.Status != corev1.DeletingStatus {
					resource.Status.Status = corev1.DeletingStatus

					return reconcile.Result{Requeue: true}, nil
				}

				if err := client.Delete(ctx); err != nil {
					logger.WithError(err).Error("trying to delete the cluster")
					resource.Status.Status = corev1.DeleteFailedStatus

					return reconcile.Result{}, err
				}

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			case "STOPPING":
				resource.Status.Status = corev1.DeletingStatus

				return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
			default:
				logger.Warn("cluster is in an unknown state, choosing to requeue instead")

				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}

		// @step: we can now delete the sysadmin token
		if err := controllers.DeleteClusterCredentialsSecret(ctx,
			t.mgr.GetClient(), resource.Namespace, resource.Name); err != nil {

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the cluster")

		return reconcile.Result{}, err
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the resource status")

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	if err := finalizer.Remove(resource); err != nil {
		logger.WithError(err).Error("removing the finalizer")

		return reconcile.Result{}, err
	}
	logger.Info("successfully deleted the cluster")

	return reconcile.Result{}, nil
}
