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

package allocations

import (
	"context"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const finalizerName = "allocations"

// Reconcile is the entrypoint for the reconciliation logic
func (a acCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"resource.name":      request.NamespacedName.Name,
		"resource.namespace": request.NamespacedName.Namespace,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// @step: retrieve the type from the api
	object := &configv1.Allocation{}
	object.Status.Status = corev1.SuccessStatus
	object.Status.Conditions = []corev1.Condition{}

	if err := a.mgr.GetClient().Get(ctx, request.NamespacedName, object); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	// @step: if we are deleting the resource, we don't need to do anything
	if finalizer.IsDeletionCandidate(object) {
		return a.Delete(ctx, object)
	}
	err := func() error {
		// @step: validate the allocation is ok
		if err := a.validateAllocation(ctx, object); err != nil {
			logger.WithError(err).Error("validating the allocation")

			return err
		}

		// @step: check the teams exist else we raise a warning
		for _, x := range object.Spec.Teams {
			if x == configv1.AllTeams {
				continue
			}
			if found, err := a.Teams().Exists(ctx, x); err != nil {
				logger.WithError(err).Error("attempting to check for the team")

				return err
			} else if !found {
				object.Status.Status = corev1.WarningStatus
				object.Status.Conditions = append(object.Status.Conditions, corev1.Condition{
					Detail:  "resource not found",
					Message: "The team " + x + " does not exist",
				})
			}
		}

		return nil
	}()
	if err != nil {
		object.Status.Status = corev1.FailureStatus
		object.Status.Conditions = []corev1.Condition{
			{
				Detail:  err.Error(),
				Message: "Failed to provision the allocation",
			},
		}
	}

	// @step we update the status of the resource
	if err := a.mgr.GetClient().Status().Update(ctx, object); err != nil {
		logger.WithError(err).Error("failed to update the resource status")

		return reconcile.Result{}, err
	}
	if err != nil {
		logger.WithError(err).Error("failed to update the resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
