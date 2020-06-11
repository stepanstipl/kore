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
	"context"
	"fmt"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "services.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service")

	// @step: retrieve the object from the api
	service := &servicesv1.Service{}
	if err := c.mgr.GetClient().Get(ctx, request.NamespacedName, service); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("failed to retrieve service from api")

		return reconcile.Result{}, err
	}
	original := service.DeepCopy()

	if service.Spec.ClusterNamespace != "" && service.Annotations[kore.AnnotationSystem] != kore.AnnotationValueTrue {
		team := service.Spec.Cluster.Namespace
		namespaceName := fmt.Sprintf("%s-%s", service.Spec.Cluster.Name, service.Spec.ClusterNamespace)
		namesplaceClaim := &clustersv1.NamespaceClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      namespaceName,
				Namespace: team,
			},
		}

		// check if the namespace specified has a NamespaceClaim
		found, err := kubernetes.GetIfExists(ctx, c.mgr.GetClient(), namesplaceClaim)
		if err != nil {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}
		if !found {
			// create NamespaceClaim
			logger.Infof("creating NamespaceClaim for namespace: %s", service.Spec.ClusterNamespace)
			namespaceClaim := &clustersv1.NamespaceClaim{
				TypeMeta: metav1.TypeMeta{
					APIVersion: clustersv1.GroupVersion.String(),
					Kind:       "NamespaceClaim",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      namespaceName,
					Namespace: team,
				},
				Spec: clustersv1.NamespaceClaimSpec{
					Name: service.Spec.ClusterNamespace,
					Cluster: corev1.Ownership{
						Group:     clustersv1.GroupVersion.Group,
						Version:   clustersv1.GroupVersion.Version,
						Kind:      "Cluster",
						Namespace: service.Spec.Cluster.Namespace,
						Name:      service.Spec.Cluster.Name,
					},
				},
			}

			if _, err := kubernetes.CreateOrUpdate(ctx, c.mgr.GetClient(), namespaceClaim); err != nil {
				logger.WithError(err).Error("trying to update or create the namespaceClaim")
				return reconcile.Result{}, err
			}
		}
	}

	spCtx := kore.NewContext(ctx, logger, c.mgr.GetClient(), c)
	provider, err := c.ServiceProviders().GetProviderForKind(spCtx, service.Spec.Kind)
	if err != nil {
		service.Status.Status = corev1.ErrorStatus
		service.Status.Message = err.Error()
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	serviceKind, err := c.ServiceKinds().Get(ctx, service.Spec.Kind)
	if err != nil {
		return reconcile.Result{}, err
	}

	servicePlan, err := c.ServicePlans().Get(ctx, service.Spec.Plan)
	if err != nil {
		return reconcile.Result{}, err
	}

	service.Status.ServiceAccessEnabled = serviceKind.Spec.ServiceAccessEnabled && !servicePlan.Spec.ServiceAccessDisabled

	finalizer := kubernetes.NewFinalizer(c.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(service) {
		return c.Delete(ctx, logger, service, finalizer, provider)
	}

	result, err := func() (reconcile.Result, error) {
		ensure := []controllers.EnsureFunc{
			c.EnsureFinalizer(logger, service, finalizer),
			c.EnsureServicePending(logger, service),
			func(ctx context.Context) (result reconcile.Result, err error) {
				return provider.Reconcile(spCtx, service)
			},
		}

		for _, handler := range ensure {
			result, err := handler(ctx)
			if err != nil {
				return reconcile.Result{}, err
			}
			if result.Requeue || result.RequeueAfter > 0 {
				return result, nil
			}
		}
		return reconcile.Result{}, nil
	}()

	if err != nil {
		logger.WithError(err).Error("failed to reconcile the service")
		if controllers.IsCriticalError(err) {
			service.Status.Status = corev1.FailureStatus
			service.Status.Message = err.Error()
		}
	}

	if err == nil && !result.Requeue && result.RequeueAfter == 0 {
		service.Status.Plan = service.Spec.Plan
		service.Status.Configuration = service.Spec.Configuration
		service.Status.Status = corev1.SuccessStatus
	}

	if err := c.mgr.GetClient().Status().Patch(ctx, service, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("failed to update the service status")

		return reconcile.Result{}, err
	}

	if err != nil {
		if controllers.IsCriticalError(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if service.Status.Status == corev1.SuccessStatus {
		return reconcile.Result{}, nil
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}
