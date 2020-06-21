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
	"time"

	"github.com/appvia/kore/pkg/utils/validation"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for removing the service
func (c Controller) Delete(
	ctx context.Context,
	logger log.FieldLogger,
	service *servicesv1.Service,
	finalizer *kubernetes.Finalizer,
	provider kore.ServiceProvider,
) (reconcile.Result, error) {
	logger.Debug("attempting to delete service from the api")

	if service.Status.Status == corev1.DeletedStatus || service.GetAnnotations()[kore.Label("finalize")] == "false" {
		if err := finalizer.Remove(service); err != nil {
			logger.WithError(err).Error("failed to remove the finalizer from the service")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if !service.Status.Status.OneOf(corev1.DeletingStatus, corev1.DeleteFailedStatus, corev1.ErrorStatus) {
		service.Status.Status = corev1.DeletingStatus
		return reconcile.Result{Requeue: true}, nil
	}

	if err := c.Teams().Team(service.Namespace).Services().CheckDelete(ctx, service); err != nil {
		if dv, ok := err.(validation.ErrDependencyViolation); ok {
			service.Status.Message = dv.Error()
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		service.Status.Status = corev1.ErrorStatus
		service.Status.Message = err.Error()
		return reconcile.Result{}, err
	}

	result, err := provider.Delete(kore.NewContext(ctx, logger, c.mgr.GetClient(), c), service)
	if err != nil {
		logger.WithError(err).Error("failed to delete the service")

		service.Status.Status = corev1.ErrorStatus
		service.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			service.Status.Status = corev1.DeleteFailedStatus
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	service.Status.Status = corev1.DeletedStatus

	// We haven't finished yet as we have to remove the finalizer in the last loop
	return reconcile.Result{Requeue: true}, nil
}
