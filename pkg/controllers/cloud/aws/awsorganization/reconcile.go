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

package awsorganization

import (
	"context"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	finalizerName = "aws-organization.kore.appvia.io"
)

// Reconcile is the entrypoint for the reconciliation logic
func (t awsCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Debug("attempting to reconcile aws organisation access")

	org := &awsv1alpha1.AWSOrganization{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, org); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}
	original := org.DeepCopy()

	// @step: create a finalizer and check if we are deleting
	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)
	if finalizer.IsDeletionCandidate(org) {
		return t.Delete(request)
	}

	// @step: ensure the defaults
	org.Status.Status = corev1.SuccessStatus
	if org.Status.Conditions == nil {
		org.Status.Conditions = &corev1.Components{}
	}

	// @step: handle the reconcile in here
	result, err := func() (reconcile.Result, error) {
		// @step: we ensure we have the credentials required
		credentials, err := EnsureCredentials(kore.NewContext(ctx, logger, t.mgr.GetClient(), t.Interface), org, org.Status.Conditions)
		if err != nil {
			return reconcile.Result{}, err
		}

		// @step: we need to ensure the credentials have the correct permission (by validating OU)
		if err := t.ValidateRoleAndOUName(ctx, org, credentials); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to reconcile the aws organization")

		org.Status.Status = corev1.FailureStatus
	}

	if err := t.mgr.GetClient().Status().Patch(ctx, org, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("updating the resource status")

		return reconcile.Result{}, nil
	}

	return result, err
}
