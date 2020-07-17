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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedConfigSpec defines the configuration for a cluster
// +k8s:openapi-gen=true
type ManagedConfigSpec struct {
	// CertificateAuthority is the location of the API certificate authority
	// +kubebuilder:validation:Required
	CertificateAuthority v1.Secret `json:"certificateAuthority,omitempty"`
	// ClientCertificate is the location of the client certificate to
	// speck back to the API
	// +kubebuilder:validation:Required
	ClientCertificate v1.Secret `json:"clientCertificate,omitempty"`
	// Domain is the domain name for this cluster
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:validation:Required
	Domain string `json:"domain,omitempty"`
}

// ManagedConfigStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type ManagedConfigStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	Conditions []corev1.Condition `json:"conditions"`
	// Phase indicates the phase of the cluster
	Phase string `json:"phase"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedConfig is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=managedconfig
type ManagedConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedConfigSpec   `json:"spec,omitempty"`
	Status ManagedConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedConfigList contains a list of Cluster
type ManagedConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedConfig `json:"items"`
}
