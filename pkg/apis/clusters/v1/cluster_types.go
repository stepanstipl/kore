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

package v1

import (
	"encoding/json"
	"fmt"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ClusterGVK is the GVK for a Cluster
var ClusterGVK = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "Cluster",
}

type ClusterComponent interface {
	runtime.Object
	corev1.StatusAware
	ApplyClusterConfiguration(cluster *Cluster) error
	ComponentDependencies() []string
}

// ClusterSpec defines the desired state of a cluster
// +k8s:openapi-gen=true
type ClusterSpec struct {
	// Kind refers to the cluster type (e.g. GKE, EKS)
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Plan is the name of the cluster plan which was used to create this cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Plan string `json:"plan"`
	// Configuration are the configuration values for this cluster
	// It will contain values from the plan + overrides by the user
	// This will provide a simple interface to calculate diffs between plan and cluster configuration
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
	// Credentials is a reference to the credentials object to use
	// +kubebuilder:validation:Required
	Credentials corev1.Ownership `json:"credentials"`
}

// ClusterStatus defines the observed state of a cluster
// +k8s:openapi-gen=true
type ClusterStatus struct {
	// APIEndpoint is the kubernetes API endpoint url
	// +kubebuilder:validation:Optional
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// CaCertificate is the base64 encoded cluster certificate
	// +kubebuilder:validation:Optional
	CaCertificate string `json:"caCertificate,omitempty"`
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
	// AuthProxyEndpoint is the endpoint of the authentication proxy for this cluster
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	AuthProxyEndpoint string `json:"authProxyEndpoint,omitempty"`
	// Status is the overall status of the cluster
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// ProviderData is provider specific data
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	ProviderData *apiextv1.JSON `json:"providerData,omitempty"`
}

// GetProviderData unmarshals the provider data into the target object
func (c *ClusterStatus) GetProviderData(v interface{}) error {
	if c.ProviderData == nil {
		return nil
	}

	if err := json.Unmarshal(c.ProviderData.Raw, v); err != nil {
		return fmt.Errorf("failed to unmarshal cluster provider data: %w", err)
	}
	return nil
}

// SetProviderData marshals the given object as provider data
func (c *ClusterStatus) SetProviderData(v interface{}) error {
	if v == nil {
		c.ProviderData = nil
		return nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster provider data: %w", err)
	}
	c.ProviderData = &apiextv1.JSON{Raw: raw}
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the plans API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=clusters
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

func (c *Cluster) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Version:   GroupVersion.Version,
		Kind:      "Cluster",
		Namespace: c.Namespace,
		Name:      c.Name,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of clusters
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}
