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
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceGroupVersionKind is the GVK for a Service
var ServiceGroupVersionKind = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "Service",
}

// ServiceSpec defines the desired state of a service
// +k8s:openapi-gen=true
type ServiceSpec struct {
	// Kind refers to the service type
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Plan is the name of the service plan which was used to create this service
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Plan string `json:"plan"`
	// Configuration are the configuration values for this service
	// It will contain values from the plan + overrides by the user
	// This will provide a simple interface to calculate diffs between plan and service configuration
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
	// Credentials is a reference to the credentials object to use
	// +kubebuilder:validation:Optional
	Credentials corev1.Ownership `json:"credentials,omitempty"`
}

// ServiceStatus defines the observed state of a service
// +k8s:openapi-gen=true
type ServiceStatus struct {
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
	// Status is the overall status of the service
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// ProviderID is the service identifier in the service provider
	// +kubebuilder:validation:Optional
	ProviderID string `json:"providerID,omitempty"`
	// ProviderData is provider specific data
	// +kubebuilder:validation:Optional
	ProviderData string `json:"providerData,omitempty"`
	// Plan is the name of the service plan which was used to create this service
	// +kubebuilder:validation:Optional
	Plan string `json:"plan,omitempty"`
	// Configuration are the applied configuration values for this service
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	Configuration apiextv1.JSON `json:"configuration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Service is a managed service instance
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=services
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

func NewService(name, namespace string) *Service {
	return &Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// Ownership creates an Ownership object
func (s *Service) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Version:   GroupVersion.Version,
		Kind:      "Service",
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceList contains a list of services
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}
