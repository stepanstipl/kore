/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// GenericSecret indicates a generic secret
	GenericSecret = "generic"
	// KubernetesSecret indicates the secrets required to speak to the api
	KubernetesSecret = "kubernetes"
)

// SecretSpec defines the desired state of Plan
// +k8s:openapi-gen=true
type SecretSpec struct {
	// Type refers to the secret type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Type string `json:"type"`
	// Description provides a summary of the secret
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// Values are the key values to the plan
	// +kubebuilder:validation:Optional
	Data map[string]string `json:"data,omitempty"`
}

// SecretStatus defines the observed state of Plan
// +k8s:openapi-gen=true
type SecretStatus struct {
	// Conditions is a set of condition which has caused an error
	// +listType=set
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status,omitempty"`
	// Verified indicates if the secret has been verified as working
	Verified *bool `json:"verified,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Secret is the Schema for the plans API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=secrets
type Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretSpec   `json:"spec,omitempty"`
	Status SecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretList contains a list of Plan
type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secret `json:"items"`
}
