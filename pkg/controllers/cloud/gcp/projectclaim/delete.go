/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package projectclaim

import (
	"context"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is the entrypoint for the reconciliation logic
func (t ctrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"project": request.NamespacedName.Name,
		"team":    request.NamespacedName.Namespace,
	})
	logger.Info("attempting to delete gcp project")

	project := &gcp.ProjectClaim{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, project); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := project.DeepCopy()

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	// @step: we need to check if the project exists and delete it
	result, err := func() (reconcile.Result, error) {
		// @step: update the resource at deleting if not done already
		if project.Status.Status != corev1.DeleteStatus {
			project.Status.Status = corev1.DeleteStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: ensure the gcp organization
		org, err := t.EnsureOrganization(ctx, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{RequeueAfter: 2 * time.Minute}, err
		}

		// @step: we need to grab the credentials from the organization and create clients
		secret, err := t.EnsureOrganizationCredentials(ctx, org, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{}, err
		}

		// @step: ensure the project is deleted if it exists
		if err := t.EnsureProjectDeleted(ctx, secret, org, project); err != nil {
			logger.WithError(err).Error("trying to ensure project deleted")

			return reconcile.Result{}, err
		}

		// @step: ensure the credentials secret has gone
		if err := t.EnsureCredentialsDeleted(ctx, project); err != nil {
			logger.WithError(err).Error("trying to ensure the project credentials are deleted")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gcp project claim")

		project.Status.Status = corev1.FailureStatus
	}

	if !result.Requeue && result.RequeueAfter <= 0 {
		// @step: we can remove the finalizer now
		if err := finalizer.Remove(project); err != nil {
			logger.WithError(err).Error("trying to remove the finalizer gcp project claim")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	// @step: update the status of the resource
	if err := t.mgr.GetClient().Status().Patch(ctx, project, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update resource status of gcp project claim")

		return reconcile.Result{}, err
	}

	return result, err
}
