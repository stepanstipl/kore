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

package podpolicy

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FilterClustersBySource returns a list of kubenetes cluster in the hub - if the
// namespace is global we retrieve all clusters, else just the local teams
func (a pspCtrl) FilterClustersBySource(ctx context.Context,
	clusters []corev1.Ownership,
	teams []string,
	namespace string) (*clustersv1.KubernetesList, error) {

	list := &clustersv1.KubernetesList{}

	// @step: is the role targetting a specific cluster
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

	if hub.IsGlobalTeam(namespace) {
		return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(""))
	}

	return list, a.mgr.GetClient().List(ctx, list, client.InNamespace(namespace))
}
