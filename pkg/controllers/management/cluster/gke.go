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

	accounts "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type gkeComponents struct {
	*Controller
}

// Components generate a graph of the things that need to be created and loaded
func (g *gkeComponents) Components(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		meta := metav1.ObjectMeta{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
		}

		switch IsAccountManaged(cluster.Spec.Credentials) {
		case true:
			v := components.Add(&gcp.ProjectClaim{ObjectMeta: meta})
			c := components.Add(&gke.GKE{ObjectMeta: meta})
			components.Edge(v, c)
		default:
			components.Add(&gke.GKE{ObjectMeta: meta})
		}

		return reconcile.Result{}, nil
	}
}

func (g *gkeComponents) Complete(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	client := g.mgr.GetClient()

	return func(ctx context.Context) (reconcile.Result, error) {
		return reconcile.Result{}, components.WalkFunc(func(v *Vertex) (bool, error) {
			switch {
			case utils.IsEqualType(v.Object, &gke.GKE{}):
				o := v.Object.(*gke.GKE)

				if err := kubernetes.PatchSpec(o, cluster.Spec.Configuration.Raw); err != nil {
					return false, err
				}

				o.Spec.Cluster = cluster.Ownership()

				switch IsAccountManaged(cluster.Spec.Credentials) {
				case true:
					o.Spec.Credentials = corev1.Ownership{
						Group:     gcp.GroupVersion.Group,
						Version:   gcp.GroupVersion.Version,
						Kind:      "ProjectClaim",
						Namespace: cluster.Namespace,
						Name:      cluster.Name,
					}
				default:
					o.Spec.Credentials = cluster.Spec.Credentials
				}

			case utils.IsEqualType(v.Object, &gcp.ProjectClaim{}):
				o := v.Object.(*gcp.ProjectClaim)

				// @step: we never touch the project claim under these circumstances
				if v.Exists {
					return true, nil
				}

				// @step: we find the matching account rule
				mgmt, err := FindAccountManagement(ctx, client, cluster.Spec.Credentials)
				if err != nil {
					return false, err
				}
				o.Spec.Organization = corev1.Ownership{
					Group:     gcp.GroupVersion.Group,
					Kind:      "Organization",
					Name:      mgmt.Spec.Organization.Name,
					Namespace: mgmt.Spec.Organization.Namespace,
					Version:   gcp.GroupVersion.Version,
				}

				switch len(mgmt.Spec.Rules) > 0 {
				case true:
					rule, found := FindAccountingRule(mgmt, cluster.Spec.Plan)
					if !found {
						return false, controllers.NewCriticalError(
							fmt.Errorf("no account rule matching plan: %q exist", cluster.Spec.Plan),
						)
					}

					// @step: we derive the account name from the rule
					name := cluster.Namespace
					if rule.Suffix != "" {
						name = fmt.Sprintf("%s-%s", name, rule.Suffix)
					}
					if rule.Prefix != "" {
						name = fmt.Sprintf("%s-%s", rule.Prefix, name)
					}

					o.Spec.ProjectName = name

				default:
					// else we are just create a project per cluster
					o.Spec.ProjectName = fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)
				}
			}

			return true, nil
		})
	}
}

// SetProviderData saves the provider data on the cluster
func (g *gkeComponents) SetProviderData(cluster *clustersv1.Cluster, components *Components) controllers.EnsureFunc {
	return func(ctx context.Context) (reconcile.Result, error) {
		return reconcile.Result{}, nil
	}
}

// IsAccountManaged checks if accounting is switched on
func IsAccountManaged(owner corev1.Ownership) bool {
	if owner.Group != accounts.GroupVersion.Group {
		return false
	}
	if owner.Version != accounts.GroupVersion.Version {
		return false
	}
	if owner.Kind != "AccountManagement" {
		return false
	}

	return true
}
