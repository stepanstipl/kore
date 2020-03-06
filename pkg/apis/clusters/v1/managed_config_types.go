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
	// +listType=set
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
