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

package awsaccount

import (
	"context"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers/cloud/aws/awsorganization"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "aws-account.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t awsCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to reconcile aws accounts")

	account := &awsv1alpha1.AWSAccount{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, account); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := account.DeepCopy()

	// @step: create a finalizer and check if we are deleting
	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(account) {
		return t.Delete(request)
	}

	// @step: ensure the defaults
	if account.Status.Conditions == nil {
		account.Status.Conditions = &corev1.Components{}
	}
	// @step: handle the reconcile in here
	result, err := func() (reconcile.Result, error) {
		// @step: we ensure we have the credentials required
		// @step: if the status is not set we should set to pending
		if account.Status.Status == "" {
			account.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: check if we need to add the finalizer
		if finalizer.NeedToAdd(account) {
			if err := finalizer.Add(account); err != nil {
				logger.WithError(err).Error("trying to add the finalizer to resource")

				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: ensure the project has access to the org
		if err := t.EnsurePermitted(ctx, account); err != nil {
			logger.WithError(err).Error("checking if account has permission to aws organization")

			return reconcile.Result{}, err
		}

		// @step: ensure the account has not been created already
		if err := t.EnsureAccountUnclaimed(ctx, account); err != nil {
			logger.WithError(err).Error("checking if account is projected")

			return reconcile.Result{}, err
		}

		// @step: ensure the aws organization
		org, err := t.EnsureOrganization(ctx, account)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the aws organization")

			return reconcile.Result{}, err
		}

		// @step: we need to grab the credentials from the organization and create clients
		korectx := kore.NewContext(ctx, logger, t.mgr.GetClient(), t.Interface)
		credentials, err := awsorganization.EnsureCredentials(korectx, org, account.Status.Conditions)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve the aws organization access")

			return reconcile.Result{}, err
		}

		// @step: ensure the aws account is created
		client, provisioned, err := t.EnsureAWSAccountProvisioned(ctx, account, org, credentials)
		if err != nil {
			logger.WithError(err).Error("trying to ensure the account is provisioned")

			return reconcile.Result{}, err
		}
		if !provisioned {

			// wait for provisioning to be complete
			return reconcile.Result{Requeue: true}, nil
		}

		// Update the status with the ID for the account now we know it's been provisioned
		account.Status.AccountID = client.GetAccountID()

		// @step: ensure initial access
		accessReady, err := t.EnsureAccessFromMasterRole(ctx, account, client)
		if err != nil {
			logger.WithError(err).Error("trying to ensure access to account is provisioned")

			return reconcile.Result{}, err
		}
		if !accessReady {
			logger.WithError(err).Error("waiting for access to be granted")

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: ensure new credentials secret is created and references updated
		if err := t.EnsureCredentials(ctx, account, client); err != nil {
			logger.WithError(err).Error("trying to ensure account credentials exist for account")

			return reconcile.Result{}, err
		}

		if err := t.EnsureCredentialsAllocation(ctx, account); err != nil {
			logger.WithError(err).Error("trying to ensure the account allocation")

			return reconcile.Result{}, err
		}
		account.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the aws account")

		account.Status.Status = corev1.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, account, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("updating the resource status")

		return reconcile.Result{}, nil
	}
	return result, nil
}
