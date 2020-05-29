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

// PlanSpec defines the desired state of Plan
// +k8s:openapi-gen=true
type PlanSpec struct {
	// DisplayName is the display name of the plan
	// +kubebuilder:validation:Optional
	DisplayName string `json:"displayName,omitempty"`
	// Summary provides a short summary for the plan
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// Description provides detailed description of this plan
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// Resource refers to the resource type this is a plan for
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Labels is a collection of labels for this plan
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// Configuration are the key+value pairs describing a cluster configuration
	// +kubebuilder:validation:Type=object
	Configuration apiextv1.JSON `json:"configuration"`
}

// PlanStatus defines the observed state of Plan
// +k8s:openapi-gen=true
type PlanStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	// +listType=set
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Plan is the Schema for the plans API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=plans
type Plan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlanSpec   `json:"spec,omitempty"`
	Status PlanStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlanList contains a list of Plan
type PlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plan `json:"items"`
}
