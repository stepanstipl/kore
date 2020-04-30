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

// ServiceProviderSpec defines the desired state of a Service provider
// +k8s:openapi-gen=true
type ServiceProviderSpec struct {
	// Type refers to the service provider type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Type string `json:"type"`
	// Description provides a summary of the provider
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Description string `json:"description"`
	// Summary provides a short title summary for the provider
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// Configuration are the key+value pairs describing a service provider
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
}

// ServiceProviderStatus defines the observed state of a service provider
// +k8s:openapi-gen=true
type ServiceProviderStatus struct {
	// Status is the overall status of the service
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceProvider is a template for a service provider
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=serviceproviders
type ServiceProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceProviderSpec   `json:"spec,omitempty"`
	Status ServiceProviderStatus `json:"status,omitempty"`
}

func NewServiceProvider(name, namespace string) *ServiceProvider {
	return &ServiceProvider{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceProvider",
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceProviderList contains a list of service providers
type ServiceProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceProvider `json:"items"`
}
