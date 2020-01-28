/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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
	// +listType
	RedirectURIs []string `json:"redirectURIs"`
}

// IDPClientStatus defines the observed state of an IDP (ID Providers)
// +k8s:openapi-gen=true
type IDPClientStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	// +listType
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
