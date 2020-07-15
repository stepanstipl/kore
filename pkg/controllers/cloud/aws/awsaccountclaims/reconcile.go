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

package awsaccountclaims

import (
	"context"
	"fmt"
	"time"

	"github.com/appvia/kore/pkg/kore"

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "aws-account-claims.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := c.logger.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile the service provider")

	// @step: retrieve the object from the api
	claim := &aws.AWSAccountClaim{}
	if err := c.mgr.GetClient().Get(ctx, request.NamespacedName, claim); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		logger.WithError(err).Error("trying to retrieve service provider from api")

		return reconcile.Result{}, err
	}
	original := claim.DeepCopy()

	// @step: are we deleting the claim?
	if !claim.GetDeletionTimestamp().IsZero() {
		return c.Delete(request)
	}

	koreCtx := kore.NewContext(ctx, logger, c.mgr.GetClient(), c)
	result, err := func() (reconcile.Result, error) {
		return controllers.DefaultEnsureHandler.Run(koreCtx,
			[]controllers.EnsureFunc{
				c.EnsureFinalizer(claim),
				c.EnsurePending(claim),
				c.EnsureAccountUnclaimed(claim),
				c.EnsureAccount(claim),
			},
		)
	}()

	if err != nil {
		logger.WithError(err).Error("failed to reconcile the gcp project claim")

		claim.Status.Status = corev1.ErrorStatus

		if controllers.IsCriticalError(err) {
			claim.Status.Status = corev1.FailureStatus
		}
	}

	if err := controllers.PatchStatus(ctx, c.mgr.GetClient(), claim, original); err != nil {
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

// EnsureFinalizer ensures the resource has the finalizer
func (c *Controller) EnsureFinalizer(claim *aws.AWSAccountClaim) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		f := kubernetes.NewFinalizer(ctx.Client(), finalizerName)
		if f.NeedToAdd(claim) {
			if err := f.Add(claim); err != nil {
				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsurePending ensures the resource is pending
func (c *Controller) EnsurePending(claim *aws.AWSAccountClaim) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if claim.Status.Status != corev1.PendingStatus {
			claim.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		return reconcile.Result{}, nil
	}
}

// EnsureAccountUnclaimed is responsible for checking the name is unique
func (c *Controller) EnsureAccountUnclaimed(claim *aws.AWSAccountClaim) controllers.EnsureFunc {
	cc := c.mgr.GetClient()

	return func(ctx kore.Context) (reconcile.Result, error) {
		// @step: ensure no claim exists outside of the team
		list := &aws.AWSAccountList{}
		if err := cc.List(ctx, list, client.InNamespace("")); err != nil {
			c.logger.WithError(err).Error("trying to retrieve all the projects")

			return reconcile.Result{}, err
		}

		// @step: ensure no other namespace is referencing this project
		for _, x := range list.Items {
			if x.Namespace == claim.Namespace {
				continue
			}

			// @note this should never happen as teams are always unique in kore
			if x.Spec.AccountName == claim.Spec.AccountName {
				return reconcile.Result{}, controllers.NewCriticalError(
					fmt.Errorf("aws account: %q is already taken by another team", claim.Spec.AccountName),
				)
			}
		}

		return reconcile.Result{}, nil
	}
}

// EnsureAccount is responsible for provisioning a project none created
func (c *Controller) EnsureAccount(claim *aws.AWSAccountClaim) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @step: check if a project exists already
		account := &aws.AWSAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      claim.Spec.AccountName,
				Namespace: claim.Namespace,
			},
		}

		found, err := kubernetes.GetIfExists(ctx, ctx.Client(), account)
		if err != nil {
			c.logger.WithError(err).Error("trying to check for account existence")

			return reconcile.Result{}, err
		}

		if !found {
			account.Spec.AccountName = claim.Spec.AccountName
			account.Spec.Organization = claim.Spec.Organization

			if _, err := kubernetes.CreateOrUpdate(ctx, ctx.Client(), account); err != nil {
				c.logger.WithError(err).Error("trying to create the account")

				return reconcile.Result{}, err
			}
			claim.Status.AccountRef = account.Ownership()

			return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
		}

		if found {
			switch account.Status.Status {
			case corev1.PendingStatus, "":
				return reconcile.Result{RequeueAfter: 15 * time.Second}, nil

			case corev1.DeletingStatus:
				return reconcile.Result{RequeueAfter: 30 * time.Second}, nil

			case corev1.FailureStatus:
				claim.Status.Status = corev1.FailureStatus
				claim.Status.Conditions = account.Status.Conditions

				return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
			}
		}

		claim.Status.Conditions = account.Status.Conditions
		claim.Status.CredentialRef = account.Status.CredentialRef
		claim.Status.AccountID = account.Status.AccountID
		claim.Status.AccountRef = account.Ownership()
		claim.Status.Status = corev1.SuccessStatus

		return reconcile.Result{}, nil
	}
}