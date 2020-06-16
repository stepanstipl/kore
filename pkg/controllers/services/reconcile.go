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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "services.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcileResult reconcile.Result, reconcileError error) {
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

	defer func() {
		if err := c.mgr.GetClient().Status().Patch(ctx, service, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("failed to update the service status")
			reconcileResult = reconcile.Result{}
			reconcileError = err
		}
	}()

	koreCtx := kore.NewContext(ctx, logger, c.mgr.GetClient(), c)

	provider, err := c.ServiceProviders().GetProviderForKind(koreCtx, service.Spec.Kind)
	if err != nil {
		if err == kore.ErrNotFound {
			service.Status.Status = corev1.ErrorStatus
			service.Status.Message = fmt.Sprintf("There is no service provider for kind %q", service.Spec.Kind)
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		service.Status.Status = corev1.ErrorStatus
		service.Status.Message = err.Error()
		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(c.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(service) {
		return c.Delete(ctx, logger, service, finalizer, provider)
	}

	result, err := func() (reconcile.Result, error) {
		ensure := []controllers.EnsureFunc{
			c.EnsureServicePending(logger, service),
			c.EnsureDependencies(logger, service),
			c.EnsureFinalizer(logger, service, finalizer),
			func(ctx context.Context) (result reconcile.Result, err error) {
				return provider.Reconcile(koreCtx, service)
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

		service.Status.Status = corev1.ErrorStatus
		service.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			service.Status.Status = corev1.FailureStatus
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	service.Status.Status = corev1.SuccessStatus
	service.Status.Plan = service.Spec.Plan
	service.Status.Configuration = service.Spec.Configuration

	return reconcile.Result{}, nil
}
