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

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedClusterRoleBindingSpec defines the desired state of Cluster role
// +k8s:openapi-gen=true
type ManagedClusterRoleBindingSpec struct {
	// Binding is the cluster role binding you wish to propagate to the clusters
	// +kubebuilder:validation:Required
	Binding rbacv1.ClusterRoleBinding `json:"binding"`
	// Clusters is used to apply the cluster role to a specific cluster
	// +kubebuilder:validation:Optional
	// +listType=set
	Clusters []corev1.Ownership `json:"clusters,omitempty"`
	// Teams is a filter on the teams
	// +kubebuilder:validation:Optional
	// +listType=set
	Teams []string `json:"teams,omitempty"`
}

// ManagedClusterRoleStatus defines the observed state of a cluster role binding
// +k8s:openapi-gen=true
type ManagedClusterRoleBindingStatus struct {
	// Conditions is a set of condition which has caused an error
	// +listType=set
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubernetes is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=managedclusterrolebinding
type ManagedClusterRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterRoleBindingSpec   `json:"spec,omitempty"`
	Status ManagedClusterRoleBindingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterRoleBindningList contains a list of Cluster
type ManagedClusterRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedClusterRoleBinding `json:"items"`
}
