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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// AllTeams is a special group name
	AllTeams = "*"
)

// AllocationSpec defines the desired state of Allocation
// +k8s:openapi-gen=true
type AllocationSpec struct {
	// Name is the name of the resource being shared
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Summary is the summary of the resource being shared
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// Resource is the resource which is being shared with another team
	// +kubebuilder:validation:Required
	Resource corev1.Ownership `json:"resource"`
	// Teams is a collection of teams the allocation is permitted to use
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Teams []string `json:"teams"`
}

// AllocationStatus defines the observed state of Allocation
// +k8s:openapi-gen=true
type AllocationStatus struct {
	// Status is the general status of the resource
	Status corev1.Status `json:"status,omitempty"`
	// Conditions is a collection of potential issues
	// +kubebuilder:validation:Optional
	Conditions []corev1.Condition `json:"conditions,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Allocation is the Schema for the allocations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=allocations,scope=Namespaced
// +kubebuilder:printcolumn:name="Summary",type="string",JSONPath=".spec.summary",description="A summary of what is being shared"
// +kubebuilder:printcolumn:name="Group",type="string",JSONPath=".spec.resource.group",description="The API group of the resource being shared"
// +kubebuilder:printcolumn:name="Resource Namespace",type="string",JSONPath=".spec.resource.namespace",description="The namespace of the resource being shared"
// +kubebuilder:printcolumn:name="Resource Name",type="string",JSONPath=".spec.resource.name",description="The name of the resource being shared"
type Allocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AllocationSpec   `json:"spec,omitempty"`
	Status AllocationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AllocationList contains a list of Allocation
type AllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Allocation `json:"items"`
}
