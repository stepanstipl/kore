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
	"encoding/json"
	"fmt"
	"strings"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"k8s.io/apimachinery/pkg/api/meta"
)

// createClusterComponents is responsible for generating the various components required
// by the cluster - i.e. the cloud provider, perhaps a VPC for EKS etc.
func createClusterComponents(c *clustersv1.Cluster) (map[string]clustersv1.ClusterComponent, error) {
	components := map[string]clustersv1.ClusterComponent{}

	provider := createProvider(c)
	components[getComponentName(provider)] = provider

	providerComponents, err := createProviderComponents(c)
	if err != nil {
		return nil, err
	}
	for _, pc := range providerComponents {
		components[getComponentName(pc)] = pc
	}

	kubernetes := clustersv1.NewKubernetes(c.Name, c.Namespace)
	kubernetes.Spec.Provider = corev1.Ownership{
		Group:     provider.GetObjectKind().GroupVersionKind().Group,
		Kind:      provider.GetObjectKind().GroupVersionKind().Kind,
		Name:      c.Name,
		Namespace: c.Namespace,
		Version:   provider.GetObjectKind().GroupVersionKind().Version,
	}
	components[getComponentName(kubernetes)] = kubernetes

	return components, nil
}

// createProvider is responsible for create the cloud provider object based on the
// kubernetes backing provider
func createProvider(c *clustersv1.Cluster) clustersv1.ClusterComponent {
	switch strings.ToLower(c.Spec.Kind) {
	case "gke":
		return gkev1alpha1.NewGKE(c.Name, c.Namespace)
	case "eks":
		return eksv1alpha1.NewEKS(c.Name, c.Namespace)
	default:
		panic(fmt.Errorf("unknown provider type: %q", c.Spec.Kind))
	}
}

// createProviderComponents generates any additionals components required by the provider -
// such as a VPC for EKS
func createProviderComponents(c *clustersv1.Cluster) ([]clustersv1.ClusterComponent, error) {
	switch strings.ToLower(c.Spec.Kind) {
	case "gke":
		return nil, nil
	case "eks":
		var components []clustersv1.ClusterComponent

		components = append(components, eksv1alpha1.NewEKSVPC(c.Name, c.Namespace))

		var config map[string]interface{}
		if err := json.Unmarshal(c.Spec.Configuration.Raw, &config); err != nil {
			return nil, err
		}
		for _, ng := range config["nodeGroups"].([]interface{}) {
			nodeGroup := ng.(map[string]interface{})
			name := c.Name + "-" + nodeGroup["name"].(string)
			components = append(components, eksv1alpha1.NewEKSNodeGroup(name, c.Namespace))
		}

		return components, nil
	default:
		panic(fmt.Errorf("unknown provider type: %q", c.Spec.Kind))
	}
}

func getClusterResourceVersion(c clustersv1.ClusterComponent) string {
	metaAccessor, _ := meta.Accessor(c)

	return metaAccessor.GetAnnotations()[kore.Label("revision")]
}

func setClusterResourceVersion(c clustersv1.ClusterComponent, resourceVersion string) {
	metaAccessor, _ := meta.Accessor(c)
	if metaAccessor.GetAnnotations() == nil {
		metaAccessor.SetAnnotations(map[string]string{})
	}

	metaAccessor.GetAnnotations()[kore.Label("revision")] = resourceVersion
}

func getComponentName(c clustersv1.ClusterComponent) string {
	meta, _ := kubernetes.GetMeta(c)

	return c.GetObjectKind().GroupVersionKind().Kind + "/" + meta.Name
}

func readyForReconcile(c clustersv1.ClusterComponent, components map[string]clustersv1.ClusterComponent) bool {
	for _, dep := range c.ComponentDependencies() {
		for _, depc := range components {
			if strings.HasPrefix(getComponentName(depc), dep) {
				status, _ := depc.GetStatus()
				if status != corev1.SuccessStatus {
					return false
				}
			}
		}
	}

	return true
}

func readyForDelete(c clustersv1.ClusterComponent, components map[string]clustersv1.ClusterComponent) bool {
	for _, comp := range components {
		if hasDependency(comp, getComponentName(c)) {
			status, _ := comp.GetStatus()
			if status != corev1.DeletedStatus {
				return false
			}
		}
	}
	return true
}

func hasDependency(c clustersv1.ClusterComponent, name string) bool {
	for _, dep := range c.ComponentDependencies() {
		if strings.HasPrefix(name, dep) {
			return true
		}
	}
	return false
}

func applyEKSVPC(eksvpc *eksv1alpha1.EKSVPC, components map[string]clustersv1.ClusterComponent) {
	for _, c := range components {
		if applier, ok := c.(eksv1alpha1.EKSVPCApplier); ok {
			applier.ApplyEKSVPC(eksvpc)
		}
	}
}
