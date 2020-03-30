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
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesSpec defines the desired state of Cluster
// +k8s:openapi-gen=true
type KubernetesSpec struct {
	// AuthProxyImage is the kube api proxy used to sso into the cluster post provision
	// +kubebuilder:validation:Optional
	AuthProxyImage string `json:"authProxyImage,omitempty"`
	// AuthProxyAllowedIPs is a list of IP address ranges (using CIDR format), which will be allowed to access the proxy
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	AuthProxyAllowedIPs []string `json:"authProxyAllowedIPs,omitempty"`
	// Cluster refers to the cluster this object belongs to
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// ClusterUsers is a collection of users from the team whom have
	// permissions across the cluster
	// +kubebuilder:validation:Optional
	// +listType=set
	ClusterUsers []ClusterUser `json:"clusterUsers,omitempty"`
	// EnableDefaultTrafficBlock indicates the cluster should default to
	// enabling blocking network policies on all namespaces
	EnableDefaultTrafficBlock *bool `json:"enableDefaultTrafficBlock,omitempty"`
	// DefaultTeamRole is role inherited by all team members
	// +kubebuilder:validation:Optional
	DefaultTeamRole string `json:"defaultTeamRole,omitempty"`
	// Domain is the domain of the cluster
	// +kubebuilder:validation:Optional
	Domain string `json:"domain,omitempty"`
	// InheritTeamMembers inherits indicates all team members are inherited
	// as having access to cluster by default.
	// +kubebuilder:validation:Optional
	InheritTeamMembers bool `json:"inheritTeamMembers,omitempty"`
	// Provider is the cloud cluster provider type for this kubernetes
	// +kubebuilder:validation:Optional
	Provider corev1.Ownership `json:"provider,omitempty"`
}

// ClusterUser defines a user and their role in the cluster
type ClusterUser struct {
	// Username is the team member the role is being applied to
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Username string `json:"username"`
	// Roles is the roles the user is permitted access to
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	Roles []string `json:"roles"`
}

// KubernetesStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type KubernetesStatus struct {
	// Endpoint is the kubernetes endpoint url
	// +kubebuilder:validation:Optional
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// CaCertificate is the base64 encoded cluster certificate
	// +kubebuilder:validation:Optional
	CaCertificate string `json:"caCertificate,omitempty"`
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
	// APIEndpoint is the endpoint of client proxy for this cluster
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Endpoint string `json:"endpoint,omitempty"`
	// Status is overall status of the workspace
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubernetes is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubernetes
type Kubernetes struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesSpec   `json:"spec,omitempty"`
	Status KubernetesStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList contains a list of Cluster
type KubernetesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kubernetes `json:"items"`
}
