/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
