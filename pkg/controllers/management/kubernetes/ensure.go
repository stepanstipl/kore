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

package kubernetes

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// EnsureResourceDeletion is responsible for cleanup the resources in the cluster
// @note: at present this is only done for EKS as GKE performs it's own cleanup
func (a k8sCtrl) EnsureResourceDeletion(ctx context.Context, object *clustersv1.Kubernetes) error {
	logger := log.WithFields(log.Fields{
		"name":      object.Name,
		"namespace": object.Namespace,
	})

	// @note: it debatable if this should be includes as the user wont see it anyhow
	object.Status.Components.SetCondition(corev1.Component{
		Name:    ComponentClusterDelete,
		Message: "attempting to clean up in-cluster kubernetes resources",
	})

	// First delete all namespaces to ensure this will work
	// @step: retrieve the provider credentials secret
	token, err := controllers.GetConfigSecret(ctx, a.mgr.GetClient(), object.Namespace, object.Name)
	if err != nil {
		if kerrors.IsNotFound(err) {
			logger.WithError(err).Warn("kubernetes secret was not found")

			object.Status.Components.SetCondition(corev1.Component{
				Name:    ComponentClusterDelete,
				Message: "Unable obtain cluster access cluster credentials",
				Detail:  err.Error(),
				Status:  corev1.FailureStatus,
			})

			return err
		}

		logger.WithError(err).Error("unable to obtain credentials secret")

		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentClusterDelete,
			Message: "Unable obtain cluster access cluster credentials",
			Detail:  err.Error(),
			Status:  corev1.FailureStatus,
		})

		return err
	}

	// @step: create a client for the remote cluster
	client, err := kubernetes.NewRuntimeClientFromConfigSecret(token)
	if err != nil {
		logger.WithError(err).Error("trying to create client from credentials secret")

		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentClusterDelete,
			Message: "Unable to access cluster using provided cluster credentials",
			Detail:  err.Error(),
			Status:  corev1.FailureStatus,
		})

		return err
	}

	if err = CleanupKoreCluster(ctx, client); err != nil {
		logger.WithError(err).Error("trying to clean up cluster resources")

		object.Status.Components.SetCondition(corev1.Component{
			Name:    ComponentClusterDelete,
			Message: "Unable to delete all cluster namespaces",
			Detail:  err.Error(),
			Status:  corev1.FailureStatus,
		})

		return err
	}

	return nil
}
