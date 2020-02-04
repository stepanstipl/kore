/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ClusterUserRolesToMap iterates the clusters users and dedups them
func ClusterUserRolesToMap(users []clustersv1.ClusterUser) map[string][]string {
	roles := make(map[string][]string)

	for _, user := range users {
		for _, role := range user.Roles {
			if list, found := roles[role]; found {
				list = append(list, user.Username)
				roles[role] = list
			} else {
				roles[role] = []string{user.Username}
			}
		}
	}

	return roles
}

// ReconcileClusterRequests builds a list of requests based on a team change
func ReconcileClusterRequests(ctx context.Context, cc client.Client, team string) ([]reconcile.Request, error) {
	logger := log.WithFields(log.Fields{
		"team": team,
	})
	logger.Info("triggering a cluster reconcilation based on upstream trigger")

	list, err := ListAllClustersInTeam(ctx, cc, team)
	if err != nil {
		logger.WithError(err).Error("trying to list teams in clusters")

		// @TODO we need way to surface these to the users
		return []reconcile.Request{}, err
	}

	return ClustersToRequests(list.Items), nil
}

// ClustersToRequests converts a collection of claims to requests
func ClustersToRequests(items []clustersv1.Kubernetes) []reconcile.Request {
	requests := make([]reconcile.Request, len(items))

	// @step: trigger the namespaceclaims to reconcile
	for i := 0; i < len(items); i++ {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      items[i].Name,
				Namespace: items[i].Namespace,
			},
		}
	}

	return requests
}

// ListAllClustersInTeam does what it says
func ListAllClustersInTeam(ctx context.Context, cc client.Client, namespace string) (*clustersv1.KubernetesList, error) {
	list := &clustersv1.KubernetesList{}

	return list, cc.List(ctx, list, client.InNamespace(namespace))

}
