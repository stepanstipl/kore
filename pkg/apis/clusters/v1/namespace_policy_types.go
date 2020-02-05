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

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamepacePolicySpec defines the desired state of Cluster role
// +k8s:openapi-gen=true
type NamepacePolicySpec struct {
	// DefaultLabels are the labels applied to all managed namespaces
	// +kubebuilder:validation:Optional
	DefaultLabels map[string]string `json:"defaultLabels,omitempty"`
	// DefaultAnnotations are default annotations applied to all managed namespaces
	// +kubebuilder:validation:Optional
	DefaultAnnotations map[string]string `json:"defaultAnnotations,omitempty"`
	// DefaultLimits are the default resource limits applied to the namespace
	// +kubebuilder:validation:Optional
	DefaultLimits *core.LimitRange `json:"defaultLimits,omitempty"`
}

// NamepacePolicyStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type NamepacePolicyStatus struct {
	// Conditions is a set of condition which has caused an error
	// +listType
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubernetes is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=namespacepolicy
type NamepacePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamepacePolicySpec   `json:"spec,omitempty"`
	Status NamepacePolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList contains a list of Managed
type NamepacePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamepacePolicy `json:"items"`
}
