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

package namespaceclaims

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileNamespaceClaims returns a list of claims to reconcile based on a change
func ReconcileNamespaceClaims(ctx context.Context, c client.Client, name, namespace string) ([]reconcile.Request, error) {
	logger := log.WithFields(log.Fields{
		"name":      name,
		"namespace": namespace,
	})
	logger.Info("triggering a namespaceclaim reconcilation based on upstream trigger")

	list, err := ListTeamNamespaceClaims(ctx, c, namespace)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of namespaceclaims in team namespace")

		// @TODO we need way to surface these to the users
		return []reconcile.Request{}, err
	}

	return NamespaceClaimsToRequests(list), nil
}

// NamespaceClaimsToRequests converts a collection of claims to requests
func NamespaceClaimsToRequests(items []clustersv1.NamespaceClaim) []reconcile.Request {
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

// ListTeamNamespaceClaims returns a list of namespaceclaims in a namespace (team in our case)
func ListTeamNamespaceClaims(ctx context.Context, c client.Client, namespace string) ([]clustersv1.NamespaceClaim, error) {
	list := &clustersv1.NamespaceClaimList{}

	if err := c.List(ctx, list, client.InNamespace(namespace)); err != nil {
		return []clustersv1.NamespaceClaim{}, err
	}

	return list.Items, nil
}
