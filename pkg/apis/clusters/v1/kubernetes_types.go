/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesSpec defines the desired state of Cluster
// +k8s:openapi-gen=true
type KubernetesSpec struct {
	// ProxyImage is the kube api proxy used to sso into the cluster post provision
	// +kubebuilder:validation:Optional
	ProxyImage string `json:"proxyImage,omitempty"`
	// ClusterUsers is a collection of users from the team whom have
	// permissions across the cluster
	// +kubebuilder:validation:Optional
	// +listType
	ClusterUsers []ClusterUser `json:"clusterUsers,omitempty"`
	// EnabledDefaultTrafficBlock indicates the cluster shoukd default to
	// enabling blocking network policies on all namespaces
	EnabledDefaultTrafficBlock *bool `json:"enabledDefaultTrafficBlock,omitempty"`
	// InheritTeamMembers inherits indicates all team members are inherited
	// as having access to cluster by default.
	// +kubebuilder:validation:Optional
	InheritTeamMembers bool `json:"inheritTeamMembers,omitempty"`
	// DefaultTeamRole is role inherited by all team members
	// +kubebuilder:validation:Optional
	DefaultTeamRole string `json:"defaultTeamRole,omitempty"`
	// Domain is the domain of the cluster
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Domain string `json:"domain"`
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
	Roles []string `json:"roles"`
}

// KubernetesStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type KubernetesStatus struct {
	// AdminToken is the kore-admin service account token which is bound to cluster-admin
	AdminToken core.Secret `json:"adminToken,omitempty"`
	// CaCertificate is the base64 encoded cluster certificate
	CaCertificate string `json:"caCertificate,omitempty"`
	// Components is a collection of component statuses
	Components corev1.Components `json:"components,omitempty"`
	// APIEndpoint is the endpoint of client proxy for this cluster
	// +kubebuilder:validation:MinLength=1
	Endpoint string `json:"endpoint,omitempty"`
	// Endpoint is the kubernetes endpoint url
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
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
