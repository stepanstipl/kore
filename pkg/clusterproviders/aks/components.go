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

package aks

import (
	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/clusterproviders"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
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

	aks := &aksv1alpha1.AKS{ObjectMeta: meta}

	components.AddComponent(&kore.ClusterComponent{
		Object:     aks,
		IsProvider: true,
	})

	kubernetesObj.Dependencies = append(kubernetesObj.Dependencies, aks)

	return nil
}

// BeforeComponentsUpdate runs after the components are loaded but before updated
// The cluster components will be provided in dependency order
func (p Provider) BeforeComponentsUpdate(ctx kore.Context, cluster *clustersv1.Cluster, components *kore.ClusterComponents) error {
	config := &Configuration{}
	if err := cluster.Spec.GetConfiguration(config); err != nil {
		return err
	}

	for _, comp := range *components {
		switch o := comp.Object.(type) {
		case *aksv1alpha1.AKS:
			o.Spec = aksv1alpha1.AKSSpec{
				AuthorizedIPRanges:      ConvertAuthorizedMasterNetworks(config.AuthorizedMasterNetworks),
				Cluster:                 cluster.Ownership(),
				Credentials:             cluster.Spec.Credentials,
				Description:             config.Description,
				DNSPrefix:               config.DNSPrefix,
				EnablePodSecurityPolicy: config.EnablePodSecurityPolicy,
				KubernetesVersion:       config.Version,
				LinuxProfile:            ConvertLinuxProfile(config.LinuxProfile),
				NetworkPlugin:           config.NetworkPlugin,
				AgentPoolProfiles:       ConvertNodePools(config.NodePools),
				EnablePrivateCluster:    config.PrivateClusterEnabled,
				Location:                config.Region,
				WindowsProfile:          ConvertWindowsProfile(config.WindowsProfile),
			}
			if config.NetworkPolicy != "" {
				o.Spec.NetworkPolicy = utils.StringPtr(config.NetworkPolicy)
			}
			clusterproviders.SetCommonTags(&o.Spec.Tags, cluster)
		}
	}

	return nil
}

// SetProviderData saves the provider data on the cluster
// The cluster components will be provided in dependency order
func (p Provider) SetProviderData(kore.Context, *clustersv1.Cluster, *kore.ClusterComponents) error {
	return nil
}
