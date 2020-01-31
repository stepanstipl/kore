/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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

package gke

import (
	"context"
	"fmt"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Delete is responsible for deleting the gke cluster
func (t *gkeCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
		"team":      request.NamespacedName.Name,
	})
	logger.Info("attempting to delete gke cluster")

	// @step: first we need to check if we have access to the credentials
	resource := &gke.GKE{}
	if err := t.mgr.GetClient().Get(ctx, request.NamespacedName, resource); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	finalizer := kubernetes.NewFinalizer(t.mgr.GetClient(), finalizerName)

	requeue, err := func() (bool, error) {
		creds, err := t.GetCredentials(ctx, resource, request.NamespacedName.Name)
		if err != nil {
			return false, err
		}

		// @step: create a cloud client for us
		client, err := NewClient(creds, resource)
		if err != nil {
			return false, err
		}

		// @step: check if the cluster exists and if so we wait or the operation or the exit
		found, err := client.Exists()
		if err != nil {
			return false, fmt.Errorf("checking if cluster exists: %s", err)
		}
		if found {
			return false, client.Delete(ctx)
		}

		return false, nil
	}()
	if err != nil {
		logger.WithError(err).Error("attempting to delete the cluster")

		return reconcile.Result{}, err
	}
	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}
	if err := finalizer.Remove(resource); err != nil {
		logger.WithError(err).Error("removing the finalizer")

		return reconcile.Result{}, err
	}
	logger.Info("successfully deleted the cluster")

	return reconcile.Result{}, nil
}
