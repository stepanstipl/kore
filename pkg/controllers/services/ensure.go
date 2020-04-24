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

	log "github.com/sirupsen/logrus"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnsureServicePending ensures the service has a pending status
func (c *Controller) EnsureServicePending(logger log.FieldLogger, service *servicesv1.Service) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
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
func (c *Controller) EnsureFinalizer(logger log.FieldLogger, service *servicesv1.Service, finalizer *kubernetes.Finalizer) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		if finalizer.NeedToAdd(service) {
			err := finalizer.Add(service)
			if err != nil {
				logger.WithError(err).Error("failed to set the finalizer")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}
		return reconcile.Result{}, nil
	}
}
