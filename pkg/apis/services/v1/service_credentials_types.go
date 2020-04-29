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

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceCredentialsSpec defines the the desired status for service credentials
// +k8s:openapi-gen=true
type ServiceCredentialsSpec struct {
	// Kind refers to the service type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Service contains the reference to the service object
	// +kubebuilder:validation:Required
	Service corev1.Ownership `json:"service,omitempty"`
	// Cluster contains the reference to the cluster where the credentials will be saved as a secret
	// +kubebuilder:validation:Required
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// ClusterNamespace is the target namespace in the cluster where the secret will be created
	// +kubebuilder:validation:Required
	ClusterNamespace string `json:"clusterNamespace,omitempty"`
	// Configuration are the configuration values for this service credentials
	// It will be used by the service provider to provision the credentials
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
}

// ServiceCredentialsStatus defines the observed state of a service
// +k8s:openapi-gen=true
type ServiceCredentialsStatus struct {
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
	// Status is the overall status of the service
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceCredentials is credentials provisioned by a service into the target namespace
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=servicecredentials
type ServiceCredentials struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceCredentialsSpec   `json:"spec,omitempty"`
	Status ServiceCredentialsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceCredentialsList contains a list of service credentials
type ServiceCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceCredentials `json:"items"`
}
