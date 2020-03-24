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

// PlanPolicySpec defines Plan JSON Schema extensions
// +k8s:openapi-gen=true
type PlanPolicySpec struct {
	// Kind refers to the cluster type this is a plan policy for
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Labels is a collection of labels for this plan policy
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// Summary provides a short title summary for the plan policy
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// Description provides a detailed description of the plan policy
	// +kubebuilder:validation:Optional
	Description string `json:"description"`
	// Properties are the
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	Properties []PlanPolicyProperty `json:"properties"`
}

// PlanPolicyStatus defines the observed state of Plan Policy
// +k8s:openapi-gen=true
type PlanPolicyStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	// +listType=set
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the plan policy
	Status corev1.Status `json:"status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlanPolicy is the Schema for the plan policies API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=planpolicies
type PlanPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlanPolicySpec   `json:"spec,omitempty"`
	Status PlanPolicyStatus `json:"status,omitempty"`
}

// PlanPolicyProperty defines a JSON schema for a given property
// +k8s:openapi-gen=true
type PlanPolicyProperty struct {
	// Name is the name of the property
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Schema is the JSON Schema definition for the given property
	// +kubebuilder:validation:Type=object
	// +kubebuilder:validation:Optional
	Schema apiextv1.JSON `json:"schema"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlanPolicyList contains a list of Plan Policies
type PlanPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlanPolicy `json:"items"`
}
