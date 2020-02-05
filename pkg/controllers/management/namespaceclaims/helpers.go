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
	"fmt"

	clusterv1 "github.com/appvia/hub-apis/pkg/apis/clusters/v1"
	orgv1 "github.com/appvia/hub-apis/pkg/apis/org/v1"
	kubev1 "github.com/appvia/kube-operator/pkg/apis/kube/v1"
	k8s "github.com/gambol99/hub-utils/pkg/kubernetes"

	"github.com/appvia/hub-apiserver/pkg/hub"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// HubNamespaceType returns a reference to the hub namespaced type
func HubNamespaceType(name string) types.NamespacedName {
	return types.NamespacedName{
		Name:      name,
		Namespace: "hub",
	}
}

// MakeTeamMembersList is responsible for retrieving the teams members
func MakeTeamMembersList(ctx context.Context, c client.Client, team string) (*orgv1.TeamMembershipList, error) {
	users := &orgv1.TeamMembershipList{}
	namespace := fmt.Sprintf("team-%s", team)

	if err := c.List(ctx, users, client.InNamespace(namespace)); err != nil {
		return users, err
	}

	return users, nil
}

// MakeClusterKubeClient is responisble for creating a kubernetes client from a cluster reference
func MakeClusterKubeClient(ctx context.Context, c client.Client, reference types.NamespacedName) (kubernetes.Interface, error) {
	// @step: retrieve the cluster resource
	cluster := &clusterv1.Kubernetes{}
	if err := c.Get(ctx, reference, cluster); err != nil {
		return nil, err
	}

	// @step: create a kubernetes client to this cluster
	return k8s.NewFromToken(cluster.Spec.Endpoint, cluster.Spec.Token, "")
}

// ReconcileNamespaceClaims returns a list of claims to reconcile based on a change
func ReconcileNamespaceClaims(ctx context.Context, c client.Client, name, namespace string) ([]reconcile.Request, error) {
	log.WithValues(
		"trigger.name", name,
		"trigger.namespace", namespace,
	).Info("triggering a namespaceclaim reconcilation based on upstream trigger")

	list, err := ListTeamNamespaceClaims(ctx, c, namespace)
	if err != nil {
		log.WithValues(
			"trigger.name", name,
			"trigger.namespace", namespace,
		).Error(err, "failed to retrieve a list of namespaceclaims in team namespace")

		// @TODO we need way to surface these to the users
		return []reconcile.Request{}, err
	}

	return NamespaceClaimsToRequests(list), nil
}

// NamespaceClaimsToRequests converts a collection of claims to requests
func NamespaceClaimsToRequests(items []kubev1.NamespaceClaim) []reconcile.Request {
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
func ListTeamNamespaceClaims(ctx context.Context, c client.Client, team string) ([]kubev1.NamespaceClaim, error) {
	list := &kubev1.NamespaceClaimList{}

	if err := c.List(ctx, list, client.InNamespace(team)); err != nil {
		return []kubev1.NamespaceClaim{}, err
	}

	return list.Items, nil
}
