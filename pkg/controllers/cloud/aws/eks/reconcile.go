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

package eks

import (
	"context"

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	core "github.com/appvia/kore/pkg/apis/core/v1"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileCredentials ensure the cluste has it's configuration
func (t eksCtrl) ReconcileCredentials(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile eks credentials")

	// @step: retrieve the resource from the api
	creds := &aws.AWSCredentials{}
	if err := t.mgr.GetClient().Get(context.Background(), request.NamespacedName, creds); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	deepcopy := creds.DeepCopy()

	creds.Status.Verified = true
	creds.Status.Status = core.SuccessStatus

	// @step: update the status of the resource
	err := t.mgr.GetClient().Status().Patch(context.Background(), creds, client.MergeFrom(deepcopy))
	if err != nil {
		logger.WithError(err).Error("trying to update the resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
