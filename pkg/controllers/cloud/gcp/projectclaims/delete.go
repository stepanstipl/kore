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

package projectclaims

import (
	"context"

	"github.com/appvia/kore/pkg/kore"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the claim
func (c *Controller) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service provider")

	// @step: retrieve the object from the api
	claim := &gcp.ProjectClaim{}
	if err := c.client.Get(ctx, request.NamespacedName, claim); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		c.logger.WithError(err).Error("trying to retrieve claim from api")

		return reconcile.Result{}, err
	}
	original := claim.DeepCopy()

	f := kubernetes.NewFinalizer(c.client, finalizerName)
	if !f.IsDeletionCandidate(claim) {
		return reconcile.Result{}, nil
	}

	koreCtx := kore.NewContext(ctx, logger, c.client, c)
	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(koreCtx,
			[]controllers.EnsureFunc{
				c.EnsureDeleting(claim),
				c.EnsureFinalizerRemoved(claim),
			},
		)
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gcp project claim")

		claim.Status.Status = corev1.ErrorStatus

		if controllers.IsCriticalError(err) {
			claim.Status.Status = corev1.DeleteFailedStatus
		}
	}

	if err := controllers.PatchStatus(ctx, c.client, claim, original); err != nil {
		logger.WithError(err).Error("trying to update the status")

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

// EnsureDeleting ensures the resource is deleting
func (c *Controller) EnsureDeleting(claim *gcp.ProjectClaim) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if claim.Status.Status != corev1.DeletingStatus {
			claim.Status.Status = corev1.DeletingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureFinalizerRemoved removes the finalizer
func (c *Controller) EnsureFinalizerRemoved(claim *gcp.ProjectClaim) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		f := kubernetes.NewFinalizer(ctx.Client(), finalizerName)

		if f.IsDeletionCandidate(claim) {
			return reconcile.Result{}, f.Remove(claim)
		}

		return reconcile.Result{}, nil
	}
}
