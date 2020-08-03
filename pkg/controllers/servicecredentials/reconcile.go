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

	"github.com/appvia/kore/pkg/controllers/helpers"

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
	finalizerName = "servicecredentials.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcileResult reconcile.Result, reconcileError error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service credentials")

	// @step: retrieve the object from the api
	creds := &servicesv1.ServiceCredentials{}
	if err := c.client.Get(ctx, request.NamespacedName, creds); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("trying to retrieve service credentials from api")

		return reconcile.Result{}, err
	}
	original := creds.DeepCopy()

	defer func() {
		if err := c.client.Status().Patch(ctx, creds, client.MergeFrom(original)); err != nil {
			if !kerrors.IsNotFound(err) {
				logger.WithError(err).Error("failed to update the service credentials status")
				reconcileResult = reconcile.Result{}
				reconcileError = err
			}
		}
	}()

	koreCtx := kore.NewContext(ctx, logger, c.client, c)

	provider, err := c.ServiceProviders().GetProviderForKind(koreCtx, creds.Spec.Kind)
	if err != nil {
		if err == kore.ErrNotFound {
			creds.Status.Status = corev1.ErrorStatus
			creds.Status.Message = fmt.Sprintf("There is no service provider for kind %q", creds.Spec.Kind)
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		creds.Status.Status = corev1.ErrorStatus
		creds.Status.Message = err.Error()
		return reconcile.Result{}, err
	}

	service, err := c.Teams().Team(creds.Spec.Service.Namespace).Services().Get(ctx, creds.Spec.Service.Name)
	if err != nil {
		if err == kore.ErrNotFound {
			creds.Status.Status = corev1.PendingStatus
			creds.Status.Message = fmt.Sprintf("Service %q does not exist", creds.Spec.Service.Name)
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(c.client, finalizerName)
	if finalizer.IsDeletionCandidate(creds) {
		return c.delete(ctx, logger, service, creds, finalizer, provider)
	}

	if !kore.IsSystemResource(creds) && !kubernetes.HasOwnerReferenceWithKind(creds, servicesv1.ServiceGVK) {
		return helpers.EnsureOwnerReference(ctx, c.client, creds, service)
	}

	if service.Status.Status != corev1.SuccessStatus {
		creds.Status.Status = corev1.PendingStatus
		creds.Status.Message = fmt.Sprintf("Service %q is not ready", creds.Spec.Service.Name)
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	result, err := func() (reconcile.Result, error) {
		ensure := []controllers.EnsureFunc{
			c.ensurePending(creds),
			c.EnsureDependencies(creds),
			c.ensureFinalizer(creds, finalizer),
			c.ensureSecret(service, creds, provider),
		}

		for _, handler := range ensure {
			result, err := handler(koreCtx)
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
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if result.Requeue || result.RequeueAfter > 0 {
		return result, nil
	}

	creds.Status.Status = corev1.SuccessStatus

	return reconcile.Result{}, nil
}
