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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecurityRule contains the definition of a security rule
// +k8s:openapi-gen=false
type SecurityRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecurityRuleSpec `json:"spec,omitempty"`
}

// SecurityRuleSpec specifies the details of a security rule
// +k8s:openapi-gen=false
type SecurityRuleSpec struct {
	// Code is the unique identifier of this rule
	Code string `json:"code,omitempty"`
	// Name is the human-readable name of this rule
	Name string `json:"name,omitempty"`
	// Description is the markdown-formatted extended description of this rule.
	Description string `json:"description,omitempty"`
	// AppliesTo is the list of resource types (e.g. Plan, Cluster) that this rule is applicable for
	AppliesTo []string `json:"appliesTo,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecurityRuleList contains a list of rules
type SecurityRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityRule `json:"items"`
}
