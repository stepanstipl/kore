/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package podpolicy

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting any pop which were created
func (a pspCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.Namespace,
		"namespace": request.Name,
	})
	logger.Debug("attempting to delete the object")

	// @step: retrieve the type from the api
	policy := &clustersv1.ManagedPodSecurityPolicy{}
	if err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, policy); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	// @step: create a finalizer and check if we are deleting
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if err := finalizer.Remove(policy); err != nil {
		log.WithError(err).Error("trying to remove the pod security policy")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
