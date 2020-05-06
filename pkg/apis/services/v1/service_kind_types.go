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
	"k8s.io/apimachinery/pkg/runtime/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceKindGVK is the GroupVersionKind for ServiceKind
var ServiceKindGVK = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "ServiceKind",
}

// ServiceKindSpec defines the state of a service kind
// +k8s:openapi-gen=true
type ServiceKindSpec struct {
	// Enabled is true if the service kind can be used
	// +kubebuilder:validation:Optional
	Enabled bool `json:"enabled"`
	// DisplayName refers to the display name of the service type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Optional
	DisplayName string `json:"displayName,omitempty"`
	// Description provides a summary of the service kind
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// Summary provides a short title summary for the service kind
	// +kubebuilder:validation:Optional
	Summary string `json:"summary,omitempty"`
	// ImageURL is a thumbnail for the service kind
	// +kubebuilder:validation:Optional
	ImageURL string `json:"imageURL,omitempty"`
	// } refers to the documentation page for this service
	// +kubebuilder:validation:Optional
	DocumentationURL string `json:"documentationURL,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceKind is a service type
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=servicekinds
type ServiceKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceKindSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceKindList contains a list of service kinds
type ServiceKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceKind `json:"items"`
}
