/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

package kubernetes

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting any bindings which were created
func (a k8sCtrl) Delete(ctx context.Context, object *clustersv1.Kubernetes) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      object.Name,
		"namespace": object.Namespace,
	})
	logger.Debug("attempting to delete the object")

	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), finalizerName)

	if err := finalizer.Remove(object); err != nil {
		log.WithError(err).Error("trying to remove the finalizer")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
