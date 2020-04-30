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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (c *Controller) delete(
	ctx context.Context,
	logger log.FieldLogger,
	serviceProvider *servicesv1.ServiceProvider,
	finalizer *kubernetes.Finalizer,
) (reconcile.Result, error) {
	logger.Debug("attempting to delete service provider from the api")

	if serviceProvider.Status.Status == corev1.DeletedStatus {
		err := finalizer.Remove(serviceProvider)
		if err != nil {
			logger.WithError(err).Error("failed to remove the finalizer from the service provider")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	original := serviceProvider.DeepCopyObject()

	serviceProvider.Status.Status = corev1.DeletingStatus

	result, err := func() (reconcile.Result, error) {
		_, err := c.ServiceProviders().Unregister(serviceProvider)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()

	if err != nil {
		logger.WithError(err).Error("failed to delete the service provider")

		serviceProvider.Status.Status = corev1.ErrorStatus
		serviceProvider.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			serviceProvider.Status.Status = corev1.DeleteFailedStatus
		}
	}

	if err == nil && !result.Requeue && result.RequeueAfter == 0 {
		serviceProvider.Status.Status = corev1.DeletedStatus
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

	// We haven't finished yet as we have to remove the finalizer in the last loop
	if serviceProvider.Status.Status == corev1.DeletedStatus {
		return reconcile.Result{Requeue: true}, nil
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, err
}
