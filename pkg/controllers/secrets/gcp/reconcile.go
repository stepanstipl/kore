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

package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils"
	gcputils "github.com/appvia/kore/pkg/utils/cloud/gcp"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// ErrMissingClientEmail indicates the client email was missing from the service account
	ErrMissingClientEmail = errors.New("client email is missing from service account key")
	// ErrMissingOrganization indicates we were unable to find the organization associated to the account
	ErrMissingOrganization = errors.New("no gcp organization found associated to service account")
	// ErrMultipleOrganizations indicates the service account as multiple orgs associated
	ErrMultipleOrganizations = errors.New("multiple gcp organizations associated to service account")
	// ErrMissingServiceAccountKey indicates the service account is missing
	ErrMissingServiceAccountKey = errors.New("secret does not have a 'key' field holding the service account")
)

// Reconcile ensures the clusters roles across all the managed clusters
func (a ctrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile generic secret roles")

	// @step: retrieve the resource from the api
	secret := &configv1.Secret{}
	if err := a.mgr.GetClient().Get(ctx, request.NamespacedName, secret); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := secret.DeepCopy()

	// @step: we only care about gcp organizational secrets
	if secret.Spec.Type != assets.GCPOrganizationalSecret.Name {
		return reconcile.Result{}, nil
	}

	// @step: we need to get if the roles exists in the service account
	result, err := func() (reconcile.Result, error) {
		// @step: set the resource to pending
		if secret.Status.Status != corev1.PendingStatus {
			secret.Status.Status = corev1.PendingStatus
			secret.Status.Conditions = []corev1.Condition{}

			return reconcile.Result{Requeue: true}, nil
		}

		// @step: decode the secret
		copied := secret.DeepCopy()
		if err := copied.Decode(); err != nil {
			return reconcile.Result{}, err
		}

		// @step: set the default and false
		secret.Status.Verified = utils.BoolPtr(false)

		// @step: check the key is set
		sa, found := copied.Spec.Data["key"]
		if !found {
			return reconcile.Result{}, ErrMissingServiceAccountKey
		}

		// @step: retrieve the client email from the service account
		account, found, err := gcputils.GetServiceAccountFromKeyFile(sa)
		if err != nil {
			return reconcile.Result{}, err
		}
		if !found {
			return reconcile.Result{}, ErrMissingClientEmail
		}

		// @step: create a resource manager client
		client, err := gcputils.CreateResourceManagerClientFromServiceAccount(sa)
		if err != nil {
			return reconcile.Result{}, err
		}

		// a list of missing roles
		var missing []string

		err = utils.Retry(ctx, 3, true, 5*time.Second, func() (bool, error) {
			list, err := gcputils.GetServiceAccountOrganizationsIDs(ctx, sa)
			if err != nil {
				logger.WithError(err).Error("trying to retrieve service account organization")

				return false, nil
			}
			if len(list) <= 0 {
				return false, ErrMissingOrganization
			}
			if len(list) > 1 {
				return false, ErrMultipleOrganizations
			}
			id := list[0]

			roles, err := gcputils.CheckOrganizationRoles(ctx, id, account, client)
			if err != nil {
				logger.WithError(err).Error("trying to check service account roles for gcp credentials")

				return false, nil
			}

			// @step: check which if any roles are missing
			for _, x := range a.RequiredRoles() {
				if !utils.Contains(x, roles) {
					missing = append(missing, x)
				}
			}

			if len(missing) > 0 {
				return false, fmt.Errorf("missing the follings: %s", strings.Join(missing, ","))
			}

			return true, nil
		})
		if err != nil {
			return reconcile.Result{}, err
		}

		secret.Status.Status = corev1.SuccessStatus
		secret.Status.Verified = utils.BoolPtr(true)
		secret.Status.Conditions = []corev1.Condition{}

		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to check the gcp organization credentials")

		secret.Status.Status = corev1.FailureStatus
		secret.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Failed trying to validate the GCP Organization secret",
		}}
	}

	if err := a.mgr.GetClient().Status().Patch(ctx, secret, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update generic secret resource status")

		return reconcile.Result{}, err
	}

	return result, nil
}

// RequiredRoles returns the roles required to provision
func (a ctrl) RequiredRoles() []string {
	return []string{
		"roles/billing.user",
		"roles/browser",
		"roles/iam.securityReviewer",
		"roles/orgpolicy.policyViewer",
		"roles/resourcemanager.projectCreator",
		"roles/resourcemanager.projectDeleter",
		"roles/viewer",
	}
}
