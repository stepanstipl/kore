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

package clusterbindings

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
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

// Reconcile ensures the clusters bindings across all the managed clusters
func (a crCtrl) Reconcile(request reconcile.Request) (reconcile.Result, error) {
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
						kore.Label("owner"): "true",
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
			Message: "Failed trying to reconcile the managed binding",
			Detail:  err.Error(),
		}}
	} else {
		binding.Status.Status = corev1.SuccessStatus
		binding.Status.Conditions = []corev1.Condition{}
	}

	if err := a.mgr.GetClient().Status().Patch(context.Background(), original, client.MergeFrom(binding)); err != nil {
		logger.WithError(err).Error("trying to update the resource status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
