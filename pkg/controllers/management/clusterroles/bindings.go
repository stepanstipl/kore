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

package clusterroles

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileBinding ensures the clusters bindings across all the managed clusters
func (a crCtrl) ReconcileBindings(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithFields(log.Fields{
		"name":      request.NamespacedName.Name,
		"namespace": request.NamespacedName.Namespace,
	})
	logger.Debug("attempting to reconcile cluster bindings")

	// @step: retrieve the resource from the api
	ctx := context.Background()
	binding := &clustersv1.ManagedClusterRoleBinding{}
	err := a.mgr.GetClient().Get(ctx, request.NamespacedName, binding)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
	original := binding.DeepCopy()

	err = func() error {
		logger.Debug("retrieving a list of applicable clusters to apply roles")

		// @step: retrieve a list of all the cluster whom are about to be syned
		list, err := a.FilterClustersBySource(ctx,
			binding.Spec.Clusters,
			binding.Spec.Teams,
			binding.Namespace,
		)
		if err != nil {
			logger.WithError(err).Error("trying to retrieve a list of clusters to change")

			return err
		}

		// @step: iterate the clusters and apply the changes
		for _, x := range list.Items {
			l := logger.WithFields(log.Fields{
				"cluster": x.Name,
				"team":    x.Namespace,
			})
			l.Debug("attempting to apply the changes to the clusters")

			key := types.NamespacedName{
				Name:      x.Name,
				Namespace: x.Namespace,
			}
			credentials := &v1.Secret{}

			// @step: retrieve the credentials for the cluster
			if err := a.mgr.GetClient().Get(ctx, key, credentials); err != nil {
				logger.WithError(err).Error("trying to retrieve cluster credentials")

				return err
			}

			// @step: create a client for the cluster
			client, err := kubernetes.NewRuntimeClientFromSecret(credentials)
			if err != nil {
				logger.WithError(err).Error("trying to create kubernetes client")

				return err
			}
			logger.Debug("attempting to apply the cluster role to remote cluster")

			if _, err := kubernetes.CreateOrUpdateManagedClusterRoleBinding(ctx, client, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: binding.Name,
					Labels: map[string]string{
						hub.Label("owner"): "true",
					},
				},
				RoleRef:  binding.Spec.Binding.RoleRef,
				Subjects: binding.Spec.Binding.Subjects,
			}); err != nil {
				logger.WithError(err).Error("trying to create or update the cluster role binding")

				return err
			}
		}

		return nil
	}()
	if err != nil {
		binding.Status.Status = corev1.FailureStatus
		binding.Status.Conditions = []corev1.Condition{{
			Message: "failed trying to reconcile the managed binding",
			Detail:  err.Error(),
		}}
	} else {
		binding.Status.Status = corev1.SuccessStatus
		binding.Status.Conditions = []corev1.Condition{}
	}

	if err := a.mgr.GetClient().Status().
		Patch(context.Background(), original, client.MergeFrom(binding)); err != nil {
		logger.WithError(err).Error("tryin to update the resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
