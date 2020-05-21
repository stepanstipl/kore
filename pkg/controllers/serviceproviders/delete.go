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

	result, err := func() (reconcile.Result, error) {
		if !serviceProvider.Status.Status.OneOf(corev1.DeletingStatus, corev1.DeleteFailedStatus, corev1.ErrorStatus) {
			serviceProvider.Status.Status = corev1.DeletingStatus
			serviceProvider.Status.Message = ""
			return reconcile.Result{Requeue: true}, nil
		}

		result, err := DeleteServices(
			kore.NewContext(ctx, logger, c.mgr.GetClient(), c),
			kore.HubAdminTeam,
			serviceProvider,
			serviceProvider.Status.Components,
		)
		if err != nil || result.Requeue || result.RequeueAfter > 0 {
			return result, err
		}

		complete, err := c.ServiceProviders().Unregister(kore.NewContext(ctx, logger, c.mgr.GetClient(), c), serviceProvider)
		if err != nil {
			return reconcile.Result{}, err
		}
		if !complete {
			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
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
		serviceProvider.Status.Message = ""
		// We haven't finished yet as we have to remove the finalizer in the last loop
		result.Requeue = true
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

	return result, nil
}
