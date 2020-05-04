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

package cluster

import (
	"context"
	"errors"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Load is responsible for loading the expected components
func (a *Controller) Load(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	client := a.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {

		return reconcile.Result{}, components.WalkFunc(func(co *Vertex) (bool, error) {
			// @step: we always ensure we have a kind on the resource
			gvk, found, err := schema.GetGroupKindVersion(co.Object)
			if err != nil || !found {
				return false, errors.New("resource gvk not found")
			}
			co.Object.GetObjectKind().SetGroupVersionKind(gvk)

			// @step: ensure we have the namespace
			SetRuntimeNamespace(co.Object, cluster.Namespace)

			// @step: check the resource exists
			co.Exists, err = kubernetes.GetIfExists(ctx, client, co.Object)
			if err != nil {
				return false, err
			}

			return true, nil
		})
	}
}

// Components is used to generate common components
func (a *Controller) Components(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		v := components.Add(&clustersv1.Kubernetes{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.Name,
				Namespace: cluster.Namespace,
			},
		})

		return reconcile.Result{}, components.WalkFunc(func(co *Vertex) (bool, error) {
			switch {
			case utils.IsEqualType(co.Object, &gke.GKE{}):
				components.Edge(co, v)
			case utils.IsEqualType(co.Object, &eks.EKSNodeGroup{}):
				components.Edge(co, v)
			}

			return true, nil
		})
	}
}

// Complete is used to complete or fill in components
func (a *Controller) Complete(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		var ownership corev1.Ownership

		return reconcile.Result{}, components.WalkFunc(func(co *Vertex) (bool, error) {
			switch {
			case utils.IsEqualType(co.Object, &clustersv1.Kubernetes{}):
				k := co.Object.(*clustersv1.Kubernetes)
				if err := kubernetes.PatchSpec(k, cluster.Spec.Configuration.Raw); err != nil {
					return false, err
				}

				k.Spec.Cluster = cluster.Ownership()
				k.Spec.Provider = ownership
			}

			return true, nil
		})
	}
}

// SetClusterStatus is responsible for ensure the status of the cluster
func (a *Controller) SetClusterStatus(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		cluster.Status.Status = corev1.PendingStatus

		// @step: walk the component and find the kubernetes type
		_ = components.WalkFunc(func(v *Vertex) (bool, error) {
			if utils.IsEqualType(v.Object, &clustersv1.Kubernetes{}) {

				k := v.Object.(*clustersv1.Kubernetes)
				cluster.Status.APIEndpoint = k.Status.APIEndpoint
				cluster.Status.AuthProxyEndpoint = k.Status.Endpoint
				cluster.Status.CaCertificate = k.Status.CaCertificate
			}

			return true, nil
		})

		if cluster.Status.Components.HasStatusForAll(corev1.SuccessStatus) {
			cluster.Status.Status = corev1.SuccessStatus
			cluster.Status.Message = "The cluster has been created successfully"
		}

		if cluster.Status.Status == corev1.SuccessStatus {
			return reconcile.Result{}, nil
		}

		if cluster.Status.Status == corev1.FailureStatus {
			return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
		}

		return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
	}
}
