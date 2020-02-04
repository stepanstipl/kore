/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package bootstrap

import (
	"context"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/finalizers"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is called when the cluster resource is being deleted
func (t bsCtrl) Delete(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})

	// @TODO need to fix this up and do somethig clever - at the moment i've not got the
	// time so we will just remove the finalizer for now

	cluster := &clusterv1.Kubernetes{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, cluster); err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	logger.Info("attempting to reconcile the deletion of kubernetes cluster")

	finalizer := finalizers.NewFinalizer(t.mgr.GetClient(), finalizerName)

	if err := finalizer.Remove(cluster); err != nil {
		logger.WithError(err).Error("failed to remove the finalizer from resource")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
