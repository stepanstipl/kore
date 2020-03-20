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

package generic

import (
	"context"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

	var verified bool

	// @step: set some default for a type
	switch secret.Spec.Type {
	case configv1.GenericSecret:
		verified = true
	default:
		// @step: ensure this is a generic secret beforehand
		if secret.Status.Verified != nil {
			return reconcile.Result{}, nil
		}
	}

	secret.Status.Conditions = []corev1.Condition{}
	secret.Status.Verified = &verified

	if err := a.mgr.GetClient().Status().Patch(ctx, secret, client.MergeFrom(original)); err != nil {
		logger.WithError(err).Error("trying to update generic secret resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
