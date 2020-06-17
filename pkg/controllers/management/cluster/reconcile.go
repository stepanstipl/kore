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
	"reflect"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	components, err := NewComponents()
	if err != nil {
		return reconcile.Result{}, err
	}

	result, err := func() (reconcile.Result, error) {
		p, err := a.Provider(cluster.Spec.Kind)
		if err != nil {
			return reconcile.Result{}, controllers.NewCriticalError(err)
		}

		return controllers.DefaultEnsureHandler.Run(ctx,
			[]controllers.EnsureFunc{
				a.AddFinalizer(cluster),
				a.SetPending(cluster),
				p.Components(cluster, components),
				a.Components(cluster, components),
				a.Load(cluster, components),
				p.Complete(cluster, components),
				a.Complete(cluster, components),
				a.Apply(cluster, components),
				p.SetProviderData(cluster, components),
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
	return func(ctx context.Context) (reconcile.Result, error) {
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
	return func(ctx context.Context) (reconcile.Result, error) {
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
func (a *Controller) Apply(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	client := a.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		var result reconcile.Result

		if cluster.Status.Components == nil {
			cluster.Status.Components = corev1.Components{}
		}

		// We walk each of the components in order, we create them if required. If the resource
		// is not yet successful we wait and requeue. If the resource has failed, we throw
		// a critical failure and stop
		err := components.WalkFunc(func(co *Vertex) (bool, error) {
			condition, found := cluster.Status.Components.GetComponent(co.String())
			if !found {
				condition = &corev1.Component{
					Name:   co.String(),
					Status: corev1.PendingStatus,
				}
			}
			condition.Resource = &corev1.Ownership{
				Group:     co.Object.GetObjectKind().GroupVersionKind().Group,
				Version:   co.Object.GetObjectKind().GroupVersionKind().Version,
				Kind:      co.Object.GetObjectKind().GroupVersionKind().Kind,
				Namespace: co.Object.(metav1.Object).GetNamespace(),
				Name:      co.Object.(metav1.Object).GetName(),
			}
			logger := a.logger.WithFields(log.Fields{
				"component": co.String(),
				"condition": condition.Status,
				"existing":  co.Exists,
			})
			logger.Debug("attempting to reconciling the component")

			if condition.Status == corev1.DeletedStatus {
				return true, nil
			}

			defer func() {
				cluster.Status.Components.SetCondition(*condition)
			}()

			// @step: the resource does not exist we can simply apply it
			if !co.Exists {
				if _, err := kubernetes.CreateOrUpdate(ctx, client, co.Object); err != nil {
					return false, err
				}
				result.Requeue = true

				return false, nil
			}

			// @step: do we need to update the resource? if the revision is different yes
			if GetClusterRevision(co.Object) != cluster.ResourceVersion {
				SetClusterRevision(co.Object, cluster.ResourceVersion)

				if _, err := kubernetes.CreateOrUpdate(ctx, client, co.Object); err != nil {
					return false, err
				}
			}

			// @check if the resource is ready to reconcile
			status, err := GetObjectStatus(co.Object)
			if err != nil {
				if err == kubernetes.ErrFieldNotFound {
					result.RequeueAfter = 30 * time.Second

					return false, nil
				}
				logger.WithError(err).Error("trying to check the component status")

				return false, err
			}

			logger.WithField(
				"status", status,
			).Debug("current state of the resource")

			switch status {
			case corev1.SuccessStatus:
				condition.Message = ""
				condition.Detail = ""
				// @try and update the status straight away
				if condition.Status != corev1.SuccessStatus {
					condition.Status = corev1.SuccessStatus

					result.Requeue = true

					return false, nil
				}
				return true, nil

			case corev1.FailureStatus:
				cluster.Status.Status = corev1.FailureStatus
				condition.Status = corev1.FailureStatus

				if cnd, err := GetObjectReasonForFailure(co.Object); err == nil {
					condition.Message = cnd.Message
					condition.Detail = cnd.Detail
				} else {
					condition.Message = "Failed to provision the resource " + co.String()
					condition.Status = corev1.FailureStatus
				}

			default:
				condition.Status = status
			}
			a.logger.Debug("waiting for resource to move to success or failed")

			result.RequeueAfter = 30 * time.Second

			return false, nil
		})

		return result, err
	}
}

// Cleanup is responsible for deleting any components no longer required
func (a *Controller) Cleanup(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	client := a.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		// @logic:
		// - we iterate the component statues
		// - we find any components which are no longer referenced and we delete
		// - we update the status of the component and we requeue
		// - if the component has been removed we can delete it

		for i := 0; i < len(cluster.Status.Components); i++ {
			required, err := IsComponentReferenced(cluster.Status.Components[i], components)
			if err != nil {
				return reconcile.Result{}, err
			}
			if required {
				continue
			}

			// @step: set the resource to deleting
			if cluster.Status.Components[i].Status != corev1.DeletingStatus {
				cluster.Status.Components[i].Status = corev1.DeletingStatus

				return reconcile.Result{Requeue: true}, nil
			}

			u := ComponentToUnstructured(cluster.Status.Components[i])

			found, err := kubernetes.GetIfExists(ctx, client, u)
			if err != nil {
				return reconcile.Result{}, err
			}
			if !found {
				cluster.Status.Components.RemoveComponent(cluster.Status.Components[i].Name)

				continue
			}

			if !IsDeleting(u) {
				if err := kubernetes.DeleteIfExists(ctx, client, u); err != nil {
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
				cluster.Status.Components[i].Status = corev1.DeleteFailedStatus
			}
		}

		return reconcile.Result{}, nil
	}
}
