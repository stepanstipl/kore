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

package servicecredentials

import (
	"context"
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/kore"
	"k8s.io/apimachinery/pkg/types"

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
	finalizerName = "servicecredentials.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service credentials")

	// @step: retrieve the object from the api
	creds := &servicesv1.ServiceCredentials{}
	if err := c.mgr.GetClient().Get(ctx, request.NamespacedName, creds); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("trying to retrieve service credentials from api")

		return reconcile.Result{}, err
	}
	original := creds.DeepCopy()

	// @step: retrieve the object from the api
	service := &servicesv1.Service{}
	if err := c.mgr.GetClient().Get(ctx, creds.Spec.Service.NamespacedName(), service); err != nil {
		logger.WithError(err).Error("failed to retrieve service from api")

		return reconcile.Result{}, err
	}

	plan := &servicesv1.ServicePlan{}
	if err := c.mgr.GetClient().Get(ctx, types.NamespacedName{Namespace: kore.HubNamespace, Name: service.Spec.Plan}, plan); err != nil {
		logger.WithError(err).Error("failed to retrieve service plan from api")

		return reconcile.Result{}, err
	}

	provider := c.ServiceProviders().GetProviderForKind(creds.Spec.Kind)
	if provider == nil {
		logger.Errorf("provider not found for service kind %q", creds.Spec.Kind)
		creds.Status.Status = corev1.FailureStatus
		creds.Status.Message = fmt.Sprintf("provider not found for service kind %q", creds.Spec.Kind)
		return reconcile.Result{}, nil
	}

	finalizer := kubernetes.NewFinalizer(c.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(creds) {
		return c.delete(ctx, logger, service, plan, creds, finalizer, provider)
	}

	result, err := func() (reconcile.Result, error) {
		ensure := []controllers.EnsureFunc{
			c.ensureFinalizer(logger, creds, finalizer),
			c.ensurePending(logger, creds),
			c.EnsureActiveCluster(logger, creds),
			c.ensureSecret(logger, service, plan, creds, provider),
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
		logger.WithError(err).Error("failed to reconcile the service credentials")

		creds.Status.Status = corev1.ErrorStatus
		creds.Status.Message = err.Error()

		if controllers.IsCriticalError(err) {
			creds.Status.Status = corev1.FailureStatus
		}
	}

	if err == nil && !result.Requeue && result.RequeueAfter == 0 {
		creds.Status.Status = corev1.SuccessStatus
	}

	if err := c.mgr.GetClient().Status().Patch(ctx, creds, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("failed to update the service credentials status")

		return reconcile.Result{}, err
	}

	if err != nil {
		if controllers.IsCriticalError(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if creds.Status.Status == corev1.SuccessStatus {
		return reconcile.Result{}, nil
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
}
