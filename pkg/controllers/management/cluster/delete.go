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

package cluster

import (
	"context"
	"errors"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the cluster and all it's resources
func (a *Controller) Delete(ctx context.Context, cluster *clustersv1.Cluster) (reconcile.Result, error) {
	a.logger.Debug("attempting to delete the cluster from the api")

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)
	if !finalizer.IsDeletionCandidate(cluster) {
		a.logger.Debug("not ready for deletion yet")

		return reconcile.Result{}, nil
	}
	original := cluster.DeepCopyObject()

	components, err := NewComponents()
	if err != nil {
		a.logger.WithError(err).Error("trying to create the components")

		return reconcile.Result{}, err
	}

	result, err := func() (reconcile.Result, error) {
		p, err := a.Provider(cluster.Spec.Kind)
		if err != nil {
			return reconcile.Result{}, controllers.NewCriticalError(err)
		}

		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				a.Deleting(cluster),
				p.Components(cluster, components),
				a.Components(cluster, components),
				a.Load(cluster, components),
				a.Remove(cluster, components),
				a.RemoveFinalizer(cluster),
			},
		)
	}()
	if err != nil {
		a.logger.WithError(err).Error("trying to delete the cluster")

		if controllers.IsCriticalError(err) {
			cluster.Status.Status = corev1.FailureStatus
			cluster.Status.Message = err.Error()
		}
	}

	if err := a.mgr.GetClient().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
		if !kerrors.IsNotFound(err) {
			a.logger.WithError(err).Error("failed to update the cluster status")

			return reconcile.Result{}, err
		}
	}

	return result, err
}

// Deleting ensures the state of the cluster is set to pending if not
func (a *Controller) Deleting(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {

		switch cluster.Status.Status {
		case corev1.SuccessStatus, corev1.FailureStatus, corev1.PendingStatus, "":
			cluster.Status.Status = corev1.DeletingStatus

			return reconcile.Result{Requeue: true}, nil

		case corev1.DeletingStatus, corev1.DeletedStatus:
			return reconcile.Result{}, nil
		}

		// else the cluster is not in a state to delete yet
		return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}
}

// Remove is responsible for removing the resources one by one
func (a *Controller) Remove(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	client := a.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		// @logic
		// - we walk the components in reverse
		// - we delete the component if required
		// - we update the status of the components in the cluster
		// - we wait for it to be delete and then to move to the next

		err := components.InverseWalkFunc(func(co *Vertex) (bool, error) {
			a.logger.WithField(
				"resource", co.String(),
			).Debug("attempting to delete the resource from the cluster")

			condition := corev1.Component{Name: co.String(), Status: corev1.DeletingStatus}

			defer func() {
				cluster.Status.Components.SetCondition(condition)
			}()

			found, err := kubernetes.CheckIfExists(ctx, client, co.Object)
			if err != nil {
				return false, err
			}
			if !found {
				condition.Status = corev1.DeletedStatus

				return true, nil
			}

			// @step: check if the resource is deleting and if not try and delete it
			condition.Status = corev1.DeletingStatus

			if !IsDeleting(co.Object) {
				if err := kubernetes.DeleteIfExists(ctx, client, co.Object); err != nil {
					return false, err
				}
			}

			status, err := GetObjectStatus(co.Object)
			if err != nil {
				return false, err
			}
			switch status {
			case corev1.DeleteFailedStatus:
				condition.Status = corev1.DeleteFailedStatus
				cluster.Status.Status = corev1.FailureStatus

				return false, controllers.NewCriticalError(errors.New("Failed trying to remove resource"))
			}

			return false, nil
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		if cluster.Status.Components.HasStatusForAll(corev1.DeletedStatus) {
			cluster.Status.Status = corev1.DeletedStatus
			cluster.Status.Message = "The cluster successfully removed all components"

			return reconcile.Result{}, nil
		}

		return reconcile.Result{RequeueAfter: 20 * time.Second}, nil
	}
}

// RemoveFinalizer is responsible for removing the finalizer
func (a *Controller) RemoveFinalizer(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	client := a.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		if cluster.Status.Status != corev1.DeletedStatus {
			return reconcile.Result{RequeueAfter: 20 * time.Second}, nil
		}

		finalizer := kubernetes.NewFinalizer(client, finalizerName)

		if err := finalizer.Remove(cluster); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}
