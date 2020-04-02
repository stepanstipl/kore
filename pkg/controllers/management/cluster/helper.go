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

	"k8s.io/apimachinery/pkg/api/meta"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
)

func createClusterComponents(c *clustersv1.Cluster) ([]clustersv1.ClusterComponent, error) {
	var components []clustersv1.ClusterComponent

	provider := createProvider(c)
	components = append(components, provider)

	createProviderComponents, err := createProviderComponents(c)
	if err != nil {
		return nil, err
	}
	components = append(components, createProviderComponents...)

	kubernetes := clustersv1.NewKubernetes(c.Name, c.Namespace)
	kubernetes.Spec.Provider = corev1.Ownership{
		Group:     provider.GetObjectKind().GroupVersionKind().Group,
		Kind:      provider.GetObjectKind().GroupVersionKind().Kind,
		Name:      c.Name,
		Namespace: c.Namespace,
		Version:   provider.GetObjectKind().GroupVersionKind().Version,
	}
	components = append(components, kubernetes)

	return components, nil
}

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

func createProviderComponents(c *clustersv1.Cluster) ([]clustersv1.ClusterComponent, error) {
	switch strings.ToLower(c.Spec.Kind) {
	case "gke":
		return nil, nil
	case "eks":
		var components []clustersv1.ClusterComponent
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
	return metaAccessor.GetLabels()[labelClusterResourceVersion]
}

func setClusterResourceVersion(c clustersv1.ClusterComponent, resourceVersion string) {
	metaAccessor, _ := meta.Accessor(c)
	if metaAccessor.GetLabels() == nil {
		metaAccessor.SetLabels(map[string]string{})
	}
	metaAccessor.GetLabels()[labelClusterResourceVersion] = resourceVersion
}
