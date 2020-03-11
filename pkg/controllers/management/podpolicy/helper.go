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

package podpolicy

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FilterClustersBySource returns a list of kubenetes cluster in the kore - if the
// namespace is global we retrieve all clusters, else just the local teams
func (a pspCtrl) FilterClustersBySource(ctx context.Context,
	clusters []corev1.Ownership,
	teams []string,
	namespace string) (*clustersv1.KubernetesList, error) {

	list := &clustersv1.KubernetesList{}

	// @step: is the role targeting a specific cluster
	if len(clusters) > 0 {
		item := &clustersv1.Kubernetes{}
		for _, x := range clusters {
			if err := a.mgr.GetClient().Get(ctx, types.NamespacedName{
				Name:      x.Name,
				Namespace: x.Namespace,
			}, item); err != nil {
				if !kerrors.IsNotFound(err) {
					return list, err
				}

				continue
			}

			list.Items = append(list.Items, *item)
		}

		return list, nil
	}

	// @step: check if it's filter down to teams
	if len(teams) > 0 {
		for _, x := range teams {
			clusters := &clustersv1.KubernetesList{}

			if err := a.mgr.GetClient().List(ctx, clusters, client.InNamespace(x)); err != nil {
				return list, err
			}
			list.Items = append(list.Items, clusters.Items...)
		}

		return list, nil
	}

	if kore.IsGlobalTeam(namespace) {
		return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(""))
	}

	return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(namespace))
}
