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
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/controllers/helpers"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "servicecredentials.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(ctx kore.Context, request reconcile.Request) (reconcileResult reconcile.Result, reconcileError error) {
	ctx.Logger().Debug("attempting to reconcile the service credentials")

	// @step: retrieve the object from the api
	creds := &servicesv1.ServiceCredentials{}
	if err := ctx.Client().Get(ctx, request.NamespacedName, creds); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		ctx.Logger().WithError(err).Error("trying to retrieve service credentials from api")

		return reconcile.Result{}, err
	}
	original := creds.DeepCopy()

	defer func() {
		if err := ctx.Client().Status().Patch(ctx, creds, client.MergeFrom(original)); err != nil {
			if !kerrors.IsNotFound(err) {
				ctx.Logger().WithError(err).Error("failed to update the service credentials status")
				reconcileResult = reconcile.Result{}
				reconcileError = err
			}
		}
	}()

	provider, err := ctx.Kore().ServiceProviders().GetProviderForKind(ctx, creds.Spec.Kind)
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

	service, err := ctx.Kore().Teams().Team(creds.Spec.Service.Namespace).Services().Get(ctx, creds.Spec.Service.Name)
	if err != nil {
		if err == kore.ErrNotFound {
			creds.Status.Status = corev1.PendingStatus
			creds.Status.Message = fmt.Sprintf("Service %q does not exist", creds.Spec.Service.Name)
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(ctx.Client(), finalizerName)
	if finalizer.IsDeletionCandidate(creds) {
		return c.delete(ctx, service, creds, finalizer, provider)
	}

	if !kore.IsSystemResource(creds) && !kubernetes.HasOwnerReferenceWithKind(creds, servicesv1.ServiceGVK) {
		return helpers.EnsureOwnerReference(ctx, ctx.Client(), creds, service)
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
		ctx.Logger().WithError(err).Error("failed to reconcile the service credentials")

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
