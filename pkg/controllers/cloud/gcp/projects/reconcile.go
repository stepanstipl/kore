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

package projects

import (
	"context"
	"errors"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "projects.gcp.compute.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t ctrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Info("attempting to reconcile gcp project")

	project := &gcp.Project{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, project); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := project.DeepCopy()

	// @step: ensure we have components in the status
	if project.Status.Conditions == nil {
		project.Status.Conditions = &corev1.Components{}
	}

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(project) {
		return t.Delete(request)
	}

	result, err := func() (reconcile.Result, error) {
		// @step: if the status is not set we should set to pending
		if project.Status.Status == "" {
			project.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: check if we need to add the finalizer
		if finalizer.NeedToAdd(project) {
			if err := finalizer.Add(project); err != nil {
				logger.WithError(err).Error("trying to add the finalizer to resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: ensure the project has access to the org
		if err := t.EnsurePermitted(ctx, project); err != nil {
			logger.WithError(err).Error("checking if project has permission to gcp organization")

			return reconcile.Result{}, err
		}

		// @step: ensure thr project has not been projected already
		if err := t.EnsureProjectUnclaimed(ctx, project); err != nil {
			logger.WithError(err).Error("checking if project is projected")

			return reconcile.Result{}, err
		}

		// @step: ensure the gcp organization
		org, err := t.EnsureOrganization(ctx, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{}, err
		}

		// @step: we need to grab the credentials from the organization and create clients
		secret, err := t.EnsureOrganizationCredentials(ctx, org, project)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the gcp organization")

			return reconcile.Result{}, err
		}

		// @step: ensure the project is created
		if err := t.EnsureProject(ctx, secret, org, project); err != nil {
			logger.WithError(err).Error("trying to ensure the project")

			return reconcile.Result{}, nil
		}

		// @step: ensure the project is linked to the billing account
		if err := t.EnsureBilling(ctx, secret, org, project); err != nil {
			logger.WithError(err).Error("trying to ensure the billing account it linked")

			return reconcile.Result{}, err
		}
		// @step: ensure the project apis are enabled
		if err := t.EnsureAPIs(ctx, secret, project); err != nil {
			logger.WithError(err).Error("trying to toggle the apis in the project")

			return reconcile.Result{}, err
		}

		// @step: ensure the service account in the project
		sa, err := t.EnsureServiceAccount(ctx, secret, project)
		if err != nil {
			logger.WithError(err).Error("trying to enable the service account in project")

			// @TODO fix when we move to the controller.EnsureFunc
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		// @guard check to ensure nothing untowards happens
		if sa == nil {
			logger.Error("the service account returned was nil")

			return reconcile.Result{}, errors.New("no service account returned")
		}

		// @step: ensure the service account key in the project
		if err := t.EnsureServiceAccountKey(ctx, secret, org, sa, project); err != nil {
			logger.WithError(err).Error("trying to ensure the service account key")

			return reconcile.Result{}, err
		}

		// @step: ensure the allocation exists
		if err := t.EnsureCredentialsAllocation(ctx, project); err != nil {
			logger.WithError(err).Error("trying to ensure the project allocation")

			return reconcile.Result{}, err
		}

		project.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gcp project")

		project.Status.Status = corev1.FailureStatus
	}

	if err := controllers.PatchStatus(ctx, t.mgr.GetClient(), project, original); err != nil {
		logger.WithError(err).Error("updating the gcp project project status")

		return reconcile.Result{}, err
	}

	return result, nil
}
