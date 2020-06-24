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

package v1beta1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// CostGVK is the GroupVersionKind for Cost
var CostGVK = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "Cost",
}

// NewCost returns a new cost
func NewCost(name, namespace string) *Cost {
	return &Cost{
		TypeMeta: metav1.TypeMeta{
			Kind:       CostGVK.Kind,
			APIVersion: GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

type CostCloudProvider string

const CostCloudProviderGoogle CostCloudProvider = "gcp"
const CostCloudProviderAmazon CostCloudProvider = "aws"
const CostCloudProviderAzure CostCloudProvider = "azure"

// CostSpec defines the desired state of cost management
// +k8s:openapi-gen=true
type CostSpec struct {
	// Enabled enables or disables the cost management feature in kore.
	// +kubebuilder:validation:Optional
	Enabled bool `json:"enabled,omitempty"`
	// InfoCredentials specifies a map of cloud provider to credentials
	// for retrieving generic pricing metadata from the providers. Adding
	// a credential for a provider here enables metadata delivery for that
	// provider.
	// +kubebuilder:validation:Optional
	InfoCredentials map[CostCloudProvider]v1.SecretReference `json:"infoCredentials,omitempty"`
	// BillingCredentials specifies a map of cloud provider to credentials for
	// retrieving specific billing details from the providers. Adding a
	// credential for a provider here enables billing integration for that provider.
	// +kubebuilder:validation:Optional
	BillingCredentials map[CostCloudProvider]v1.SecretReference `json:"billingCredentials,omitempty"`
}

// CostStatus defines the observed state of a cost management service
// +k8s:openapi-gen=true
type CostStatus struct {
	// Status is the overall status of the service
	// +kubebuilder:validation:Optional
	Status corev1.Status `json:"status,omitempty"`
	// Message is the description of the current status
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
	// Components is a collection of component statuses
	// +kubebuilder:validation:Optional
	Components corev1.Components `json:"components,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cost is the Schema for the cost API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=costs
type Cost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CostSpec   `json:"spec,omitempty"`
	Status CostStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CostList contains a list of Cost
type CostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cost `json:"items"`
}

type Region struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Continent struct {
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}

type ContinentList struct {
	Items []Continent `json:"items"`
}

type PriceType string

const OnDemand PriceType = "OnDemand"
const Spot PriceType = "Spot"
const PreEmptible PriceType = "PreEmptible"

type InstanceType struct {
	Category string                `json:"category"`
	Type     string                `json:"type"`
	Prices   map[PriceType]float64 `json:"prices"`
	Cpus     float64               `json:"cpusPerVm"`
	Mem      float64               `json:"memPerVm"`
}

type InstanceTypeList struct {
	Items []InstanceType `json:"items"`
}
