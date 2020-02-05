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

package assets

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDefaultClusterRoles defines a collection of cluster roles
func GetDefaultClusterRoles() []clustersv1.ManagedClusterRole {
	return []clustersv1.ManagedClusterRole{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kore-nsadmin",
			},
			Spec: clustersv1.ManagedClusterRoleSpec{
				Enabled: true,
				Rules: []rbacv1.PolicyRule{
					{
						NonResourceURLs: []string{
							"/swagger*",
							"/swaggerapi",
							"/swaggerapi/*",
							"/version",
						},
						Verbs: []string{"get"},
					},
					{
						APIGroups: []string{
							"apps",
							"batch",
							"extensions",
							"networking.k8s.io",
						},
						Resources: []string{
							"cronjobs",
							"deployments",
							"deployments/rollback",
							"deployments/scale",
							"ingresses",
							"jobs",
							"networkpolicies",
							"replicasets",
							"replicasets/scale",
							"replicationcontrollers/scale",
							"statefulsets",
							"statefulsets/scale",
						},
						Verbs: []string{"*"},
					},
					{
						APIGroups: []string{
							"policy",
						},
						Resources: []string{
							"poddisruptionbudgets",
						},
						Verbs: []string{"*"},
					},
					{
						APIGroups: []string{""},
						Resources: []string{
							"configmaps",
							"endpoints",
							"persistentvolumeclaims",
							"persistentvolumes",
							"pods",
							"pods/attach",
							"pods/exec",
							"pods/log",
							"pods/portforward",
							"secrets",
							"serviceaccounts",
							"services",
						},
						Verbs: []string{"*"},
					},
					{
						APIGroups: []string{"autoscaling"},
						Resources: []string{"*"},
						Verbs:     []string{"*"},
					},
					{
						APIGroups: []string{"*"},
						Resources: []string{"*"},
						Verbs: []string{
							"get",
							"watch",
							"list",
						},
					},
					{
						APIGroups: []string{
							"certmanager.k8s.io",
						},
						Resources: []string{
							"certificates",
							"challenges",
							"orders",
						},
						Verbs: []string{"*"},
					},
				},
			},
		},
	}
}
