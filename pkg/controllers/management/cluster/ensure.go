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
	"fmt"

	"github.com/appvia/kore/pkg/kore"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	"github.com/appvia/kore/pkg/serviceproviders/application"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Load is responsible for loading the expected components
func (a *Controller) Load(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		for _, comp := range *components {
			comp.Object.SetNamespace(cluster.Namespace)

			if err := comp.Load(ctx); err != nil {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}
}

func (a *Controller) setComponents(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		kubernetesObj := &clustersv1.Kubernetes{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.Name,
				Namespace: cluster.Namespace,
			},
		}

		components.Add(kubernetesObj)

		kubeAppManager, err := a.createService(ctx, cluster, "kube-app-manager")
		if err != nil {
			return reconcile.Result{}, err
		}
		components.Add(kubeAppManager, kubernetesObj)

		fluxHelmOperator, err := a.createService(ctx, cluster, "flux-helm-operator")
		if err != nil {
			return reconcile.Result{}, err
		}

		components.Add(fluxHelmOperator, kubeAppManager)

		return reconcile.Result{}, nil
	}
}

func (a *Controller) setProviderComponents(provider kore.ClusterProvider, cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		if err := provider.SetComponents(ctx, cluster, components); err != nil {
			return reconcile.Result{}, err
		}

		if err := components.Sort(); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}
}

func (a *Controller) createService(ctx context.Context, cluster *clustersv1.Cluster, name string) (*servicesv1.Service, error) {
	servicePlan, err := a.ServicePlans().Get(ctx, "app-"+name)
	if err != nil {
		return nil, fmt.Errorf("failed to get service plan %q: %w", "app-"+name, err)
	}
	service := application.CreateSystemServiceFromPlan(
		*servicePlan,
		corev1.MustGetOwnershipFromObject(cluster),
		cluster.Name+"-"+name,
		cluster.Namespace,
	)
	return &service, nil
}

func (a *Controller) beforeComponentsUpdate(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		providerComponent := components.Find(func(comp kore.ClusterComponent) bool { return comp.IsProvider })

		for _, comp := range *components {
			switch o := comp.Object.(type) {
			case *clustersv1.Kubernetes:
				if err := kubernetes.PatchSpec(o, cluster.Spec.Configuration.Raw); err != nil {
					return reconcile.Result{}, err
				}

				o.Spec.Cluster = cluster.Ownership()
				if providerComponent != nil {
					o.Spec.Provider = corev1.MustGetOwnershipFromObject(providerComponent.Object)
				}
			}
		}

		return reconcile.Result{}, nil
	}
}

// SetClusterStatus is responsible for ensure the status of the cluster
func (a *Controller) SetClusterStatus(cluster *clustersv1.Cluster, components *kore.ClusterComponents) controllers.EnsureFunc {
	return func(ctx kore.Context) (reconcile.Result, error) {
		// @step: walk the component and find the kubernetes type
		for _, comp := range *components {
			switch o := comp.Object.(type) {
			case *clustersv1.Kubernetes:
				cluster.Status.APIEndpoint = o.Status.APIEndpoint
				cluster.Status.AuthProxyEndpoint = o.Status.Endpoint
				cluster.Status.CaCertificate = o.Status.CaCertificate
			}
		}

		cluster.Status.Status = corev1.SuccessStatus
		cluster.Status.Message = "The cluster has been created successfully"

		return reconcile.Result{}, nil
	}
}
