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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ ClusterProviderComponents = &koreComponents{}

// koreComponents is for reconciling the kore cluster itself
type koreComponents struct {
	*Controller
}

// Components adds any components that are required for the kore admin cluster
func (e *koreComponents) Components(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {

		// Example:

		// v := components.Add(&clustersv1.Kubernetes{
		// 	ObjectMeta: metav1.ObjectMeta{
		// 		Name:      cluster.Name,
		// 		Namespace: cluster.Namespace,
		// 	},
		// })

		// cloudInfo, err := e.createService(ctx, cluster, "cloudinfo")
		// if err != nil {
		// 	return reconcile.Result{}, err
		// }
		// cloudInfoV := components.Add(cloudInfo)
		// components.Edge(v, cloudInfoV)

		return reconcile.Result{}, nil
	}
}

// CompleteClusterComponents is used to fill in the resources if required
func (e *koreComponents) Complete(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		return reconcile.Result{}, nil
	}
}

func (e *koreComponents) SetProviderData(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		return reconcile.Result{}, nil
	}
}
