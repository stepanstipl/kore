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
	// SystemManaged indicates the secret is managed by kore and cannot be changed
	SystemManaged *bool `json:"systemManaged,omitempty"`
	// Conditions is a set of condition which has caused an error
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
