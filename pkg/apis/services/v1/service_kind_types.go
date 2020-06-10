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

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
	// ServiceAccessEnabled is true if the service provider can create service access for this service kind
	// +kubebuilder:validation:Optional
	ServiceAccessEnabled bool `json:"serviceAccessEnabled"`
	// DisplayName refers to the display name of the service type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Optional
	DisplayName string `json:"displayName,omitempty"`
	// Summary provides a short title summary for the service kind
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// Description is a detailed description of the service kind
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// ImageURL is a thumbnail for the service kind
	// +kubebuilder:validation:Optional
	ImageURL string `json:"imageURL,omitempty"`
	// DocumentationURL refers to the documentation page for this service
	// +kubebuilder:validation:Optional
	DocumentationURL string `json:"documentationURL,omitempty"`
	// Schema is the JSON schema for the plan
	// +kubebuilder:validation:Optional
	Schema string `json:"schema,omitempty"`
	// CredentialSchema is the JSON schema for credentials created for service using this plan
	// +kubebuilder:validation:Optional
	CredentialSchema string `json:"credentialSchema,omitempty"`
	// ProviderData is provider specific data
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	ProviderData *apiextv1.JSON `json:"providerData,omitempty"`
}

func (s *ServiceKindSpec) GetProviderData(v interface{}) error {
	if s.ProviderData == nil {
		return nil
	}

	if err := json.Unmarshal(s.ProviderData.Raw, v); err != nil {
		return fmt.Errorf("failed to unmarshal service provider data: %w", err)
	}
	return nil
}

func (s *ServiceKindSpec) SetProviderData(v interface{}) error {
	if v == nil {
		s.ProviderData = nil
		return nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal service kind provider data: %w", err)
	}
	s.ProviderData = &apiextv1.JSON{Raw: raw}
	return nil
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
