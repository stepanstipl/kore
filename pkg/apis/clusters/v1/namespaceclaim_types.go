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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceClaimSpec defines the desired state of NamespaceClaim
// +k8s:openapi-gen=true
type NamespaceClaimSpec struct {
	// Cluster is the cluster the namespace resides
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster"`
	// Name is the name of the namespace to create
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=3
	Name string `json:"name"`
	// Annotations is a series of annotations on the namespace
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a series of labels for the namespace
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// Limits are the limits placs on the namespace
	// +kubebuilder:validation:Optional
	Limits v1.LimitRange `json:"limits,omitempty"`
}

// NamespaceClaimStatus defines the observed state of NamespaceClaim
// +k8s:openapi-gen=true
type NamespaceClaimStatus struct {
	// Status is the status of the namespace
	Status corev1.Status `json:"status"`
	// Conditions is a series of things that caused the failure if any
	// +listType
	Conditions []corev1.Condition `json:"conditions"`
	// Phase is used to hold the current phase of the resource
	Phase string `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceClaim is the Schema for the namespaceclaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=namespaceclaims,scope=Namespaced
type NamespaceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceClaimSpec   `json:"spec,omitempty"`
	Status NamespaceClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceClaimList contains a list of NamespaceClaim
type NamespaceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceClaim `json:"items"`
}
