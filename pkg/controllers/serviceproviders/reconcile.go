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

package serviceproviders

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "serviceprovider.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service provider")

	// @step: retrieve the object from the api
	serviceProvider := &servicesv1.ServiceProvider{}
	if err := c.mgr.GetClient().Get(ctx, request.NamespacedName, serviceProvider); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("trying to retrieve service provider from api")

		return reconcile.Result{}, err
	}
	original := serviceProvider.DeepCopy()

	finalizer := kubernetes.NewFinalizer(c.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(serviceProvider) {
		return c.delete(ctx, logger, serviceProvider, finalizer)
	}

	result, err := func() (reconcile.Result, error) {
		ensure := []controllers.EnsureFunc{
			c.ensureFinalizer(serviceProvider, finalizer),
			c.ensurePending(serviceProvider),
			func(ctx context.Context) (result reconcile.Result, err error) {
				provider, err := c.ServiceProviders().Register(serviceProvider)
				if err != nil {
					return reconcile.Result{}, err
				}

				for _, plan := range provider.Plans() {
					plan.Namespace = kore.HubNamespace
					exists, err := kubernetes.CheckIfExists(ctx, c.mgr.GetClient(), &plan)
					if err != nil {
						return reconcile.Result{}, err
					}

					if !exists {
						if err := c.mgr.GetClient().Create(ctx, &plan); err != nil {
							return reconcile.Result{}, err
						}
					}
				}

				return reconcile.Result{}, nil
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
		logger.WithError(err).Error("failed to reconcile the service provider")

		serviceProvider.Status.Status = corev1.ErrorStatus
		serviceProvider.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			serviceProvider.Status.Status = corev1.FailureStatus
		}
	}

	if err == nil && !result.Requeue && result.RequeueAfter == 0 {
		serviceProvider.Status.Status = corev1.SuccessStatus
	}

	if err := c.mgr.GetClient().Status().Patch(ctx, serviceProvider, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("failed to update the service provider status")

		return reconcile.Result{}, err
	}

	if err != nil {
		if controllers.IsCriticalError(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if serviceProvider.Status.Status == corev1.SuccessStatus {
		return reconcile.Result{}, nil
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
}
