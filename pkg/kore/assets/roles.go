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
