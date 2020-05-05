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

package gkecredentials

import (
	"context"
	"fmt"
	"strings"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	gcputils "github.com/appvia/kore/pkg/utils/cloud/gcp"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t gkeCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to reconcile gke credentials")

	resource := &gke.GKECredentials{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := resource.DeepCopy()

	result, err := func() (reconcile.Result, error) {
		var verified bool
		resource.Status.Conditions = []corev1.Condition{}
		resource.Status.Verified = &verified

		// @step: set the status to pending if none set
		if resource.Status.Status == "" {
			resource.Status.Status = corev1.PendingStatus

			return reconcile.Result{Requeue: true}, nil
		}

		// for backwards-compatibility, use the key (Account) set on the GKECredentials resource, if it exists
		var key string = resource.Spec.Account
		if key == "" {
			// @step: we need to grab the secret
			secret, err := controllers.GetDecodedSecret(ctx, t.mgr.GetClient(), resource.Spec.CredentialsRef)
			if err != nil {
				return reconcile.Result{}, err
			}
			key = secret.Spec.Data["service_account_key"]
		}

		// @step: create the client to verify the permissions
		permitted, missing, err := gcputils.CheckServiceAccountPermissions(ctx,
			resource.Spec.Project,
			key,
			requiredPermissions(),
		)
		if err != nil {
			logger.WithError(err).Error("trying to verify the gke credentials")

			return reconcile.Result{}, err
		}
		resource.Status.Verified = &permitted
		resource.Status.Status = corev1.SuccessStatus

		if !permitted {
			logger.Error("gke credentials has not passed validation")

			return reconcile.Result{}, fmt.Errorf("service account is missing: %s permissions", strings.Join(missing, ","))
		}

		return reconcile.Result{}, err
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the gke credentials")

		resource.Status.Status = corev1.FailureStatus
		resource.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "Either the credentials are invalid or we've encountered an error verifying",
		}}
	}

	// @step: update the status of the resource
	if err := t.mgr.GetClient().Status().Patch(ctx, resource, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update the gke credentials resource status")

		return reconcile.Result{}, err
	}

	return result, err
}

// requirePermissions is a list of permissions required for gke credentials
func requiredPermissions() []string {
	return []string{
		"container.clusterRoleBindings.create",
		"container.clusterRoleBindings.get",
		"container.clusterRoles.bind",
		"container.clusterRoles.create",
		"container.clusters.create",
		"container.clusters.delete",
		"container.clusters.getCredentials",
		"container.clusters.list",
		"container.operations.get",
		"container.operations.list",
		"container.podSecurityPolicies.create",
		"container.secrets.get",
		"container.serviceAccounts.create",
		"container.serviceAccounts.get",
		"iam.serviceAccounts.actAs",
		"iam.serviceAccounts.get",
		"iam.serviceAccounts.list",
		"resourcemanager.projects.get",
	}
}
