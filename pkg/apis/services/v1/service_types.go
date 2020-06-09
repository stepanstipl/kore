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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceGroupVersionKind is the GroupVersionKind for Service
var ServiceGVK = schema.GroupVersionKind{
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
	// Cluster contains the reference to the cluster where the service will be created
	// +kubebuilder:validation:Optional
	Cluster corev1.Ownership `json:"cluster,omitempty"`
	// ClusterNamespace is the target namespace in the cluster where there the service will be created
	// +kubebuilder:validation:Optional
	ClusterNamespace string `json:"clusterNamespace,omitempty"`
	// Configuration are the configuration values for this service
	// It will contain values from the plan + overrides by the user
	// This will provide a simple interface to calculate diffs between plan and service configuration
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	Configuration *apiextv1.JSON `json:"configuration,omitempty"`
	// ConfigurationFrom is a way to load configuration values from alternative sources, e.g. from secrets
	// The values from these sources will override any existing keys defined in Configuration
	// +kubebuilder:validation:Optional
	// +listType=set
	ConfigurationFrom []corev1.ConfigurationFromSource `json:"configurationFrom,omitempty"`
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
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	ProviderData *apiextv1.JSON `json:"providerData,omitempty"`
	// Plan is the name of the service plan which was used to create this service
	// +kubebuilder:validation:Optional
	Plan string `json:"plan,omitempty"`
	// Configuration are the applied configuration values for this service
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	Configuration *apiextv1.JSON `json:"configuration,omitempty"`
}

func (s *ServiceStatus) GetProviderData(v interface{}) error {
	if s.ProviderData == nil {
		return nil
	}

	if err := json.Unmarshal(s.ProviderData.Raw, v); err != nil {
		return fmt.Errorf("failed to unmarshal service provider data: %w", err)
	}
	return nil
}

func (s *ServiceStatus) SetProviderData(v interface{}) error {
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
			Kind:       ServiceGVK.Kind,
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// Ownership creates an Ownership object
func (s Service) Ownership() corev1.Ownership {
	return corev1.Ownership{
		Group:     GroupVersion.Group,
		Version:   GroupVersion.Version,
		Kind:      "Service",
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

// NeedsUpdate returns true if the plan or the configuration changed compared to the status
func (s Service) NeedsUpdate() bool {
	if s.Spec.Plan != s.Status.Plan {
		return true
	}

	var raw1, raw2 []byte

	if s.Spec.Configuration != nil {
		raw1 = s.Spec.Configuration.Raw
	}
	if s.Status.Configuration != nil {
		raw2 = s.Status.Configuration.Raw
	}

	return !bytes.Equal(raw1, raw2)
}

func (s Service) GetConfiguration() *apiextv1.JSON {
	return s.Spec.Configuration
}

func (s Service) GetConfigurationFrom() []corev1.ConfigurationFromSource {
	return s.Spec.ConfigurationFrom
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceList contains a list of services
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}

// PriorityServiceSlice is used for sorting services by the priority annotation
// +k8s:openapi-gen=false
// +kubebuilder:object:generate=false
// +k8s:deepcopy-gen=false
type PriorityServiceSlice []Service

func (p PriorityServiceSlice) Len() int {
	return len(p)
}

func (p PriorityServiceSlice) Less(i, j int) bool {
	prioi, _ := strconv.Atoi(p[i].Annotations["kore.appvia.io/priority"])
	prioj, _ := strconv.Atoi(p[j].Annotations["kore.appvia.io/priority"])
	return prioi != 0 && (prioj == 0 || prioi < prioj)
}

func (p PriorityServiceSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
