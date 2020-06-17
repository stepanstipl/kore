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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Rule contains the definition of a alert rule
type Rule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RuleSpec `json:"spec,omitempty"`
}

// RuleSpec specifies the details of a alert rule
// +k8s:openapi-gen=true
type RuleSpec struct {
	// RuleID is a unique identifier for this rule
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	RuleID string `json:"ruleID,omitempty"`
	// Severity is the importance of the rule
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Severity string `json:"severity"`
	// Source is the provider of the rule i.e. prometheus, or a named source
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Source string `json:"source"`
	// Summary is a summary of the rule
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
	// RawRule is the underlying rule definition
	// +kubebuilder:validation:Required
	RawRule string `json:"rawRule"`
	// Resource is the resource the alert is for
	// +kubebuilder:validation:Required
	Resource corev1.Ownership `json:"resource"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuleList contains a list of rules
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rule `json:"items"`
}
