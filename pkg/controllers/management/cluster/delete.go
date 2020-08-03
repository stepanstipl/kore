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
	"errors"
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/kore"

	"github.com/appvia/kore/pkg/utils/validation"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the cluster and all it's resources
func (c *Controller) Delete(ctx kore.Context, cluster *clustersv1.Cluster) (reconcile.Result, error) {
	ctx.Logger().Debug("attempting to delete the cluster from the api")

	original := cluster.DeepCopyObject()

	result, err := func() (reconcile.Result, error) {
		provider, exists := kore.GetClusterProvider(cluster.Spec.Kind)
		if !exists {
			return reconcile.Result{}, controllers.NewCriticalError(fmt.Errorf("%q cluster provider is invalid", cluster.Spec.Kind))
		}

		components := &kore.ClusterComponents{}

		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				c.Deleting(cluster),
				c.CheckDelete(cluster),
				c.setComponents(cluster, components),
				c.setProviderComponents(provider, cluster, components),
				c.Load(cluster, components),
				c.Cleanup(cluster, components),
				c.Remove(cluster, components),
				c.RemoveFinalizer(cluster),
			},
		)
	}()
	if err != nil {
		ctx.Logger().WithError(err).Error("trying to delete the cluster")

		if controllers.IsCriticalError(err) {
			cluster.Status.Status = corev1.FailureStatus
			cluster.Status.Message = err.Error()
		}
	}

	if err := ctx.Client().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
		if !kerrors.IsNotFound(err) {
			ctx.Logger().WithError(err).Error("failed to update the cluster status")

			return reconcile.Result{}, err
		}
	}

	return result, err
}

// Deleting ensures the state of the cluster is set to pending if not
func (c *Controller) Deleting(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		cluster.Status.Message = ""

		switch cluster.Status.Status {
		case corev1.SuccessStatus, corev1.PendingStatus, "":
			cluster.Status.Status = corev1.DeletingStatus
			return reconcile.Result{Requeue: true}, nil
		default:
			return reconcile.Result{}, nil
		}
	}
}

// CheckDelete checks whether the cluster can be deleted
func (c *Controller) CheckDelete(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if err := ctx.Kore().Teams().Team(cluster.Namespace).Clusters().CheckDelete(ctx, cluster); err != nil {
			if dv, ok := err.(validation.ErrDependencyViolation); ok {
				cluster.Status.Message = dv.Error()
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
}

// Remove is responsible for removing the resources one by one
func (c *Controller) Remove(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @logic
		// - we walk the components in reverse
		// - we delete the component if required
		// - we update the status of the components in the cluster
		// - we wait for it to be delete and then to move to the next

		for i := len(*components) - 1; i >= 0; i-- {
			result, err := c.removeComponent(ctx, cluster, components, (*components)[i])
			if err != nil || result.Requeue || result.RequeueAfter > 0 {
				return result, err
			}
		}

		cluster.Status.Status = corev1.DeletedStatus
		cluster.Status.Message = "The cluster successfully removed all components"

		return reconcile.Result{}, nil
	}
}

func (c *Controller) removeComponent(
	ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents, comp *kore.ClusterComponent,
) (reconcile.Result, error) {
	statusComp, exists := cluster.Status.Components.GetComponent(comp.ComponentName())
	if !exists || statusComp.Status == corev1.DeletedStatus {
		return reconcile.Result{}, nil
	}

	ctx.Logger().WithField(
		"resource", comp.ComponentName(),
	).Debug("attempting to delete the resource from the cluster")

	condition := corev1.Component{Name: comp.ComponentName(), Status: corev1.DeletingStatus}

	defer func() {
		cluster.Status.Components.SetCondition(condition)
	}()

	found, err := kubernetes.CheckIfExists(ctx, ctx.Client(), comp.Object)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !found {
		if comp.AfterDelete != nil {
			if err := comp.AfterDelete(ctx, cluster, comp, components); err != nil {
				return reconcile.Result{}, err
			}
		}

		condition.Status = corev1.DeletedStatus

		return reconcile.Result{}, nil
	}

	// @step: check if the resource is deleting and if not try and delete it
	condition.Status = corev1.DeletingStatus

	if !IsDeleting(comp.Object) {
		if err := kubernetes.DeleteIfExists(ctx, ctx.Client(), comp.Object); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	status, err := GetObjectStatus(comp.Object)
	if err != nil {
		return reconcile.Result{}, err
	}
	switch status {
	case corev1.DeleteFailedStatus:
		condition.Status = corev1.DeleteFailedStatus
		cluster.Status.Status = corev1.FailureStatus

		return reconcile.Result{}, controllers.NewCriticalError(errors.New("failed trying to remove resource"))
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// RemoveFinalizer is responsible for removing the finalizer
func (c *Controller) RemoveFinalizer(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if cluster.Status.Status != corev1.DeletedStatus {
			return reconcile.Result{RequeueAfter: 20 * time.Second}, nil
		}

		finalizer := kubernetes.NewFinalizer(ctx.Client(), finalizerName)

		if err := finalizer.Remove(cluster); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}
