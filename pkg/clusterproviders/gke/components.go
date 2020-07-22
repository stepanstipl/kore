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

package gke

import (
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcpv1alpha1 "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/clusterproviders"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetComponents adds all provider-specific cluster components and updates dependencies if required
func (p Provider) SetComponents(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	kubernetesObj := components.Find(func(comp kore.ClusterComponent) bool {
		_, ok := comp.Object.(*clustersv1.Kubernetes)
		return ok
	})

	meta := metav1.ObjectMeta{
		Name:      cluster.Name,
		Namespace: cluster.Namespace,
	}

	gke := &gkev1alpha1.GKE{ObjectMeta: meta}
	var gkeDependencies []kubernetes.Object

	if clusterproviders.IsAccountManaged(cluster.Spec.Credentials) {
		projectClaim := &gcpv1alpha1.ProjectClaim{ObjectMeta: meta}
		components.Add(projectClaim)
		gkeDependencies = []kubernetes.Object{projectClaim}
	}

	components.AddComponent(&kore.ClusterComponent{
		Object:       gke,
		Dependencies: gkeDependencies,
		IsProvider:   true,
	})

	kubernetesObj.Dependencies = append(kubernetesObj.Dependencies, gke)

	return nil
}

// BeforeComponentsUpdate runs after the components are loaded but before updated
// The cluster components will be provided in dependency order
func (p Provider) BeforeComponentsUpdate(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	for _, comp := range *components {
		switch o := comp.Object.(type) {
		case *gkev1alpha1.GKE:
			if err := kubernetes.PatchSpec(o, cluster.Spec.Configuration.Raw); err != nil {
				return err
			}

			o.Spec.Cluster = cluster.Ownership()
			clusterproviders.CheckCommonTags(&o.Spec.Tags, cluster)

			switch clusterproviders.IsAccountManaged(cluster.Spec.Credentials) {
			case true:
				o.Spec.Credentials = corev1.Ownership{
					Group:     gcpv1alpha1.GroupVersion.Group,
					Version:   gcpv1alpha1.GroupVersion.Version,
					Kind:      "ProjectClaim",
					Namespace: cluster.Namespace,
					Name:      cluster.Name,
				}
			default:
				o.Spec.Credentials = cluster.Spec.Credentials
			}

		case *gcpv1alpha1.ProjectClaim:
			// @step: we never touch the project claim under these circumstances
			if comp.Exists() {
				continue
			}

			// @step: we find the matching account rule
			mgmt, err := clusterproviders.FindAccountManagement(ctx, cluster.Spec.Credentials)
			if err != nil {
				return err
			}
			o.Spec.Organization = corev1.Ownership{
				Group:     gcpv1alpha1.GroupVersion.Group,
				Kind:      "Organization",
				Name:      mgmt.Spec.Organization.Name,
				Namespace: mgmt.Spec.Organization.Namespace,
				Version:   gcpv1alpha1.GroupVersion.Version,
			}

			switch len(mgmt.Spec.Rules) > 0 {
			case true:
				rule, found := clusterproviders.FindAccountingRule(mgmt, cluster.Spec.Plan)
				if !found {
					return controllers.NewCriticalError(
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
	}

	return nil
}

// SetProviderData saves the provider data on the cluster
// The cluster components will be provided in dependency order
func (p Provider) SetProviderData(kore.Context, *clustersv1.Cluster, *kore.ClusterComponents) error {
	return nil
}
