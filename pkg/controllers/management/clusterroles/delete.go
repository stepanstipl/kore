/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package clusterroles

import (
	"context"
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Reconcile removes the roles
func (a crCtrl) Delete(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to delete managed cluster roles")

	// @step: retrieve the resource from the api
	role := &clustersv1.ManagedClusterRole{}
	err := a.mgr.GetClient().Get(context.Background(), request.NamespacedName, role)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := role.DeepCopy()
	finalizer := kubernetes.NewFinalizer(a.mgr.GetClient(), rolesFinalizer)

	// @step: grab the clusters
	// @step: retrieve a list of cluster which this role applies
	list, err := a.FilterClustersBySource(context.Background(),
		role.Spec.Clusters,
		role.Spec.Teams,
		role.Namespace)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of clusters")

		role.Status.Status = corev1.FailureStatus
		role.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed trying to retrieve list of clusters to apply",
		}}

		if err := a.mgr.GetClient().Status().Patch(context.Background(), role, client.MergeFrom(original)); err != nil {
			logger.WithError(err).Error("trying to update the resource")

			return reconcile.Result{}, err
		}

		return reconcile.Result{}, err
	}

	err = func() error {
		// @step: we iterate the clusters and apply the roles
		for _, cluster := range list.Items {
			logger := logger.WithFields(log.Fields{
				"cluster": cluster.Name,
				"team":    cluster.Namespace,
			})
			logger.Debug("attempting to remove the managed role in cluster")

			// @step: retrieve the credentials for the cluster
			credentials := &v1.Secret{}
			if err := a.mgr.GetClient().Get(context.Background(), types.NamespacedName{
				Name:      cluster.Name,
				Namespace: cluster.Namespace,
			}, credentials); err != nil {
				logger.WithError(err).Error("trying to retrieve cluster credentials")

				return err
			}

			// @step: create a client for the cluster
			client, err := kubernetes.NewRuntimeClientFromSecret(credentials)
			if err != nil {
				logger.WithError(err).Error("trying to create kubernetes client")

				return err
			}
			logger.Debug("creatin the managed cluster role in cluster")

			// @step: update or create the role
			name := fmt.Sprintf("kore:managed:%s", role.Name)
			if err := kubernetes.DeleteClusterRoleIfExists(context.Background(), client, name); err != nil {
				logger.WithError(err).Error("trying to remove the managed role")

				return err
			}
		}

		return nil
	}()
	if err != nil {
		logger.WithError(err).Error("trying to remove the managed cluster role")

		role.Status.Status = corev1.FailureStatus
		role.Status.Conditions = []corev1.Condition{{
			Detail:  err.Error(),
			Message: "failed trying to remove managed role from one of more clusters",
		}}

		return reconcile.Result{}, err
	}

	// @step: remove the finalizer
	if err := finalizer.Remove(role); err != nil {
		logger.WithError(err).Error("trying to remove the finalizer from resource")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
