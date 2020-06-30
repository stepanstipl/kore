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
	"fmt"
	"reflect"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "cluster.clusters.kore.appvia.io"
)

var (
	// ClusterRevisionName is the annotation name
	ClusterRevisionName = kore.Label("clusterRevision")
)

// Reconcile is the entrypoint for the reconciliation logic
func (a *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := a.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the cluster")

	// @step: retrieve the object from the api
	cluster := &clustersv1.Cluster{}
	if err := a.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("trying to retrieve cluster from api")

		return reconcile.Result{}, err
	}
	original := cluster.DeepCopy()

	if cluster.Annotations[kore.AnnotationSystem] == kore.AnnotationValueTrue {
		cluster.Status.Status = corev1.SuccessStatus
		if err := a.mgr.GetClient().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("failed to update the cluster status")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(cluster) {
		return a.Delete(ctx, cluster)
	}

	// @logic:
	// - we retrieve the cloud provider (this is responsible for knowing which components to create)
	// - we generate on each iteration what we need
	// - we try and load the components if they exist
	// - we then use the provider to fill in any components from other i.e. a eks <- vpc
	// - we apply the components in order when theirs dependents are ready

	result, err := func() (reconcile.Result, error) {
		provider, exists := kore.GetClusterProvider(cluster.Spec.Kind)
		if !exists {
			return reconcile.Result{}, controllers.NewCriticalError(fmt.Errorf("%q cluster provider is invalid", cluster.Spec.Kind))
		}

		components := &kore.ClusterComponents{}

		koreCtx := kore.NewContext(ctx, logger, a.mgr.GetClient(), a)
		return controllers.DefaultEnsureHandler.Run(koreCtx,
			[]controllers.EnsureFunc{
				a.AddFinalizer(cluster),
				a.SetPending(cluster),
				a.setComponents(cluster, components),
				a.setProviderComponents(provider, cluster, components),
				a.Load(cluster, components),
				func(ctx kore.Context) (reconcile.Result, error) {
					return reconcile.Result{}, provider.BeforeComponentsUpdate(ctx, cluster, components)
				},
				a.beforeComponentsUpdate(cluster, components),
				a.Apply(cluster, components),
				func(ctx kore.Context) (reconcile.Result, error) {
					return reconcile.Result{}, provider.SetProviderData(ctx, cluster, components)
				},
				a.Cleanup(cluster, components),
				a.SetClusterStatus(cluster, components),
			},
		)
	}()

	if err != nil {
		logger.WithError(err).Error("trying to ensure the cluster")

		if controllers.IsCriticalError(err) {
			cluster.Status.Status = corev1.FailureStatus
			cluster.Status.Message = err.Error()
		}
	}

	if !reflect.DeepEqual(cluster, original) {
		if err := a.mgr.GetClient().Status().Patch(ctx, cluster, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to patch the cluster status")

			return reconcile.Result{}, err
		}
	}

	return result, err
}

// AddFinalizer ensures the finalizer is on the resource
func (a *Controller) AddFinalizer(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)
		if finalizer.NeedToAdd(cluster) {
			if err := finalizer.Add(cluster); err != nil {
				a.logger.WithError(err).Error("trying to add the finalizer")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// SetPending ensures the state of the cluster is set to pending if not
func (a *Controller) SetPending(cluster *clustersv1.Cluster) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		switch cluster.Status.Status {
		case corev1.DeletingStatus:
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}

		if cluster.Status.Status == "" {
			cluster.Status.Status = corev1.PendingStatus
			return reconcile.Result{Requeue: true}, nil
		}

		cluster.Status.Status = corev1.PendingStatus
		cluster.Status.Message = ""

		return reconcile.Result{}, nil
	}
}

// Apply is responsible for applying the component and updating the component status
func (a *Controller) Apply(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if cluster.Status.Components == nil {
			cluster.Status.Components = corev1.Components{}
		}

		// We walk each of the components in order, we create them if required. If the resource
		// is not yet successful we wait and requeue. If the resource has failed, we throw
		// a critical failure and stop
		for _, comp := range *components {
			result, err := a.applyComponent(ctx, cluster, comp)
			if err != nil || result.Requeue || result.RequeueAfter > 0 {
				return result, err
			}
		}

		return reconcile.Result{}, nil
	}
}

func (a *Controller) applyComponent(ctx context.Context, cluster *clustersv1.Cluster, comp *kore.ClusterComponent) (reconcile.Result, error) {
	condition, found := cluster.Status.Components.GetComponent(comp.ComponentName())
	if !found {
		condition = &corev1.Component{
			Name:   comp.ComponentName(),
			Status: corev1.PendingStatus,
		}
	}
	defer func() {
		cluster.Status.Components.SetCondition(*condition)
	}()

	ownership := corev1.MustGetOwnershipFromObject(comp.Object)
	condition.Resource = &ownership
	logger := a.logger.WithFields(log.Fields{
		"component": comp.ComponentName(),
		"condition": condition.Status,
		"existing":  comp.Exists,
	})
	logger.Debug("attempting to reconciling the component")

	if !comp.Exists || GetClusterRevision(comp.Object) != cluster.ResourceVersion {
		SetClusterRevision(comp.Object, cluster.ResourceVersion)

		annotations := comp.Object.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}
		annotations[kore.AnnotationSystem] = kore.AnnotationValueTrue
		annotations[kore.AnnotationReadOnly] = kore.AnnotationValueTrue

		comp.Object.SetAnnotations(annotations)

		if _, err := kubernetes.CreateOrUpdate(ctx, a.mgr.GetClient(), comp.Object); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{Requeue: true}, nil
	}

	// @check if the resource is ready to reconcile
	status, err := GetObjectStatus(comp.Object)
	if err != nil {
		if err == kubernetes.ErrFieldNotFound {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}

		logger.WithError(err).Error("trying to check the component status")
		return reconcile.Result{}, err
	}

	logger.WithField(
		"status", status,
	).Debug("current state of the resource")

	condition.Status = status
	condition.Message = ""
	condition.Detail = ""

	switch status {
	case corev1.SuccessStatus:
		return reconcile.Result{}, nil
	case corev1.FailureStatus:
		if cnd, err := GetObjectReasonForFailure(comp.Object); err == nil {
			condition.Message = cnd.Message
			condition.Detail = cnd.Detail
		} else {
			condition.Message = "Failed to provision the resource"
		}
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// Cleanup is responsible for deleting any components no longer required
func (a *Controller) Cleanup(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @logic:
		// - we iterate the component statues
		// - we find any components which are no longer referenced and we delete
		// - we update the status of the component and we requeue
		// - if the component has been removed we can delete it

		for i := 0; i < len(cluster.Status.Components); i++ {
			statusComponent := cluster.Status.Components[i]
			if statusComponent.Resource == nil {
				continue
			}

			comp := components.Find(func(comp kore.ClusterComponent) bool {
				return kore.IsOwner(comp.Object, *statusComponent.Resource)
			})
			if comp != nil {
				continue
			}

			// @step: set the resource to deleting
			if statusComponent.Status != corev1.DeletingStatus {
				statusComponent.Status = corev1.DeletingStatus
			}

			u := ComponentToUnstructured(statusComponent)

			found, err := kubernetes.GetIfExists(ctx, a.mgr.GetClient(), u)
			if err != nil {
				return reconcile.Result{}, err
			}
			if !found {
				cluster.Status.Components.RemoveComponent(statusComponent.Name)

				continue
			}

			if !IsDeleting(u) {
				if err := kubernetes.DeleteIfExists(ctx, a.mgr.GetClient(), u); err != nil {
					return reconcile.Result{}, err
				}
			}

			// @check if the resource is ready to reconcile
			status, err := GetObjectStatus(u)
			if err != nil {
				return reconcile.Result{}, err
			}

			if status == corev1.DeleteFailedStatus {
				cluster.Status.Status = corev1.FailureStatus
				statusComponent.Status = corev1.DeleteFailedStatus
			}

			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}

		return reconcile.Result{}, nil
	}
}
