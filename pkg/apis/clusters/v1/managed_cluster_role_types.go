/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedClusterRoleSpec defines the desired state of Cluster role
// +k8s:openapi-gen=true
type ManagedClusterRoleSpec struct {
	// Clusters is used to apply to one of more clusters role to a specific cluster
	// +kubebuilder:validation:Optional
	// +listType=set
	Clusters []corev1.Ownership `json:"clusters,omitempty"`
	// Description provides a short summary of the nature of the role
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=10
	Description string `json:"description,omitempty"`
	// Enabled indicates if the role is enabled or not
	// +kubebuilder:validation:Optional
	Enabled bool `json:"enabled,omitempty"`
	// Rules are the permissions on the role
	// +kubebuilder:validation:Required
	// +listType=set
	Rules []rbacv1.PolicyRule `json:"rules,omitempty"`
	// Teams is used to filter the clusters to apply by team references
	// +kubebuilder:validation:Optional
	// +listType=set
	Teams []string `json:"teams,omitempty"`
}

// ManagedClusterRoleStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type ManagedClusterRoleStatus struct {
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
// +kubebuilder:resource:path=managedclusterrole
type ManagedClusterRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterRoleSpec   `json:"spec,omitempty"`
	Status ManagedClusterRoleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList contains a list of Cluster
type ManagedClusterRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedClusterRole `json:"items"`
}
