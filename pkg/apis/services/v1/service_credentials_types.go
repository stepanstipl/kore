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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceCredentialsGVK is the GroupVersionKind for ServiceCredentials
var ServiceCredentialsGVK = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "ServiceCredentials",
}

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
	// SecretName is the Kubernetes Secret's name that will contain the service access information
	// If not set the secret's name will default to `Name`
	// +kubebuilder:validation:Optional
	SecretName string `json:"secretName,omitempty"`
	// Configuration are the configuration values for this service credentials
	// It will be used by the service provider to provision the credentials
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	Configuration *apiextv1.JSON `json:"configuration,omitempty"`
}

func (s *ServiceCredentialsSpec) GetConfiguration(v interface{}) error {
	if s.Configuration == nil {
		return nil
	}

	if err := json.Unmarshal(s.Configuration.Raw, v); err != nil {
		return fmt.Errorf("failed to unmarshal service credentials configuration: %w", err)
	}
	return nil
}

func (s *ServiceCredentialsSpec) SetConfiguration(v interface{}) error {
	if v == nil {
		s.Configuration = nil
		return nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal service credentials configuration: %w", err)
	}
	s.Configuration = &apiextv1.JSON{Raw: raw}
	return nil
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
	// ProviderID is the service credentials identifier in the service provider
	// +kubebuilder:validation:Optional
	ProviderID string `json:"providerID,omitempty"`
	// ProviderData is provider specific data
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	ProviderData *apiextv1.JSON `json:"providerData,omitempty"`
}

func (s *ServiceCredentialsStatus) GetProviderData(v interface{}) error {
	if s.ProviderData == nil {
		return nil
	}

	if err := json.Unmarshal(s.ProviderData.Raw, v); err != nil {
		return fmt.Errorf("failed to unmarshal service provider data: %w", err)
	}
	return nil
}

func (s *ServiceCredentialsStatus) SetProviderData(v interface{}) error {
	if v == nil {
		s.ProviderData = nil
		return nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal service provider data: %w", err)
	}
	s.ProviderData = &apiextv1.JSON{Raw: raw}
	return nil
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

func NewServiceCredentials(name, namespace string) *ServiceCredentials {
	return &ServiceCredentials{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceCredentialsGVK.Kind,
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// SecretName returns with the Kubernetes Secret's name which will contain the service credential details
func (s ServiceCredentials) SecretName() string {
	if s.Spec.SecretName != "" {
		return s.Spec.SecretName
	}

	return s.Name
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceCredentialsList contains a list of service credentials
type ServiceCredentialsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceCredentials `json:"items"`
}
