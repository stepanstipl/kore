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

package features

import (
	"context"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/controllers/helpers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (c *Controller) delete(
	ctx context.Context,
	logger log.FieldLogger,
	feature *configv1.KoreFeature,
	finalizer *kubernetes.Finalizer,
) (reconcile.Result, error) {
	logger.Debug("attempting to delete feature from the api")

	if feature.Status.Status == corev1.DeletedStatus {
		err := finalizer.Remove(feature)
		if err != nil {
			logger.WithError(err).Error("failed to remove the finalizer from the service provider")
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	if !feature.Status.Status.OneOf(corev1.DeletingStatus, corev1.DeleteFailedStatus, corev1.ErrorStatus) {
		feature.Status.Status = corev1.DeletingStatus
		feature.Status.Message = ""
		return reconcile.Result{Requeue: true}, nil
	}

	if err := c.kore.Features().CheckDelete(ctx, feature); err != nil {
		if dv, ok := err.(validation.ErrDependencyViolation); ok {
			feature.Status.Message = dv.Error()
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		feature.Status.Status = corev1.ErrorStatus
		feature.Status.Message = err.Error()
		return reconcile.Result{}, err
	}

	result, err := func() (reconcile.Result, error) {
		result, err := helpers.DeleteServices(
			kore.NewContext(ctx, logger, c.client, c.kore),
			kore.HubAdminTeam,
			feature,
			&feature.Status.Components,
		)
		if err != nil || result.Requeue || result.RequeueAfter > 0 {
			return result, err
		}

		return reconcile.Result{}, nil
	}()

	if err != nil {
		logger.WithError(err).Error("failed to delete the feature")

		feature.Status.Status = corev1.ErrorStatus
		feature.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			feature.Status.Status = corev1.DeleteFailedStatus
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	feature.Status.Status = corev1.DeletedStatus
	feature.Status.Message = ""

	// We haven't finished yet as we have to remove the finalizer in the last loop
	return reconcile.Result{Requeue: true}, nil
}
