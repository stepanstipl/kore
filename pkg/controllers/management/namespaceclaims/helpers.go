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

	list, err := ListTeamNamespaceClaims(ctx, c, name)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve a list of namespaceclaims in team namespace")

		// @TODO we need way to surface these to the users
		return []reconcile.Request{}, err
	}

	logger.WithField(
		"namespaceclaims", len(list),
	).Debug("triggering a namespaceclaim reconcilation based on upstream trigger")

	return ToRequests(list), nil
}

// ToRequests converts a collection of claims to requests
func ToRequests(items []clustersv1.NamespaceClaim) []reconcile.Request {
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
