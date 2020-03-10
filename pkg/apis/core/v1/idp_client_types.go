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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// IDPClientSpec defines the spec for a IDP client
// +k8s:openapi-gen=true
type IDPClientSpec struct {
	// DisplayName
	// +kubebuilder:validation:Required
	DisplayName string `json:"displayName"`
	// Secret for OIDC client
	// +kubebuilder:validation:Required
	Secret string `json:"secret"`
	// ID of OIDC client
	// +kubebuilder:validation:Required
	ID string `json:"id"`
	// RedirectURIs where to send client after IDP auth
	// +kubebuilder:validation:Required
	// +listType=set
	RedirectURIs []string `json:"redirectURIs"`
}

// IDPClientStatus defines the observed state of an IDP (ID Providers)
// +k8s:openapi-gen=true
type IDPClientStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	// +listType=set
	Conditions []Condition `json:"conditions"`
	// Status is overall status of the IDP configuration
	Status Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDPClient is the Schema for the class API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=oidclient,scope=Namespaced
type IDPClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IDPClientSpec   `json:"spec,omitempty"`
	Status IDPClientStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDPClientList contains a list of IDP clients
type IDPClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDPClient `json:"items"`
}
