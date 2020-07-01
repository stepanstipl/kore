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

package services

import (
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/controllers/helpers"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/appvia/kore/pkg/kore"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureServicePending ensures the service has a pending status
func (c *Controller) EnsureServicePending(service *servicesv1.Service) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if service.Status.Status == "" {
			service.Status.Status = corev1.PendingStatus
			return reconcile.Result{Requeue: true}, nil
		}

		if service.Status.Status != corev1.PendingStatus {
			service.Status.Status = corev1.PendingStatus
		}
		return reconcile.Result{}, nil
	}
}

// EnsureFinalizer ensures the service has a finalizer
func (c *Controller) EnsureFinalizer(service *servicesv1.Service, finalizer *kubernetes.Finalizer) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if finalizer.NeedToAdd(service) {
			err := finalizer.Add(service)
			if err != nil {
				ctx.Logger().WithError(err).Error("failed to set the finalizer")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}
		return reconcile.Result{}, nil
	}
}

func (c *Controller) EnsureDependencies(service *servicesv1.Service) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		serviceKind, err := c.ServiceKinds().Get(ctx, service.Spec.Kind)
		if err != nil {
			if err == kore.ErrNotFound {
				service.Status.Status = corev1.PendingStatus
				service.Status.Message = fmt.Sprintf("Service kind %q does not exist", service.Spec.Kind)
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
			return reconcile.Result{}, err
		}

		servicePlan, err := c.ServicePlans().Get(ctx, service.Spec.Plan)
		if err != nil {
			if err == kore.ErrNotFound {
				service.Status.Status = corev1.PendingStatus
				service.Status.Message = fmt.Sprintf("Service plan %q does not exist", service.Spec.Plan)
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
			return reconcile.Result{}, err
		}

		service.Status.ServiceAccessEnabled = serviceKind.Spec.ServiceAccessEnabled && !servicePlan.Spec.ServiceAccessDisabled

		if !kore.IsSystemResource(service) && !kubernetes.HasOwnerReferenceWithKind(service, clustersv1.ClusterGVK) {
			cluster, err := c.Teams().Team(service.Spec.Cluster.Namespace).Clusters().Get(ctx, service.Spec.Cluster.Name)
			if err != nil {
				if err == kore.ErrNotFound {
					service.Status.Status = corev1.PendingStatus
					service.Status.Message = fmt.Sprintf("Cluster %q does not exist", service.Spec.Cluster.Name)
					return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
				}
				return reconcile.Result{}, err
			}

			return helpers.EnsureOwnerReference(ctx, ctx.Client(), service, cluster)
		}

		if service.Spec.ClusterNamespace != "" && !kore.IsSystemResource(service) && !kubernetes.HasOwnerReferenceWithKind(service, clustersv1.NamespaceClaimGVK) {
			name := fmt.Sprintf("%s-%s", service.Spec.Cluster.Name, service.Spec.ClusterNamespace)

			namespaceClaim, err := c.Teams().Team(service.Namespace).NamespaceClaims().Get(ctx, name)
			if err != nil {
				if kerrors.IsNotFound(err) || err == kore.ErrNotFound {
					service.Status.Status = corev1.PendingStatus
					service.Status.Message = fmt.Sprintf("Namespace claim does not exist for namespace %q", service.Spec.ClusterNamespace)
					return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
				}
				return reconcile.Result{}, err
			}

			return helpers.EnsureOwnerReference(ctx, ctx.Client(), service, namespaceClaim)
		}

		return reconcile.Result{}, nil
	}
}
