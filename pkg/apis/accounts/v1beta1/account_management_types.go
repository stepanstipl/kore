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

const (
	// AccountManagementName is the name of the CRD
	AccountManagementName = "Account"
)

// AccountManagementSpec defines the desired state of accounting for a provider
// I've a feeling this will probably need provider specific attributes are some point
// +k8s:openapi-gen=true
type AccountManagementSpec struct {
	// Provider is the name of provider which maps to the cluster kind
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`
	// Rules is a set of rules for this provider
	// +kubebuilder:validation:Optional
	// +listType=set
	Rules []*AccountsRule `json:"rules,omitempty"`
	// Organization is the underlying organizational resource (only require if more than one)
	// +kubebuilder:validation:Required
	Organization corev1.Ownership `json:"organization,omitempty"`
}

// AccountsRule defines a rule for the provider
type AccountsRule struct {
	// Name is the given name of the rule
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Description provides an optional description for the account rule
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// Plans is a list of plans permitted
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +listType=set
	Plans []string `json:"plans"`
	// Exact override any values in prefix and suffix uses that for the account name
	// +kubebuilder:validation:Optional
	Exact string `json:"exact,omitempty"`
	// Suffix is the applied suffix
	// +kubebuilder:validation:Optional
	Suffix string `json:"suffix,omitempty"`
	// Prefix is a prefix for the account name
	// +kubebuilder:validation:Optional
	Prefix string `json:"prefix,omitempty"`
	// Labels a collection of labels to apply the account
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
}

// AccountManagementStatus defines the observed state of Allocation
// +k8s:openapi-gen=true
type AccountManagementStatus struct {
	// Status is the general status of the resource
	Status corev1.Status `json:"status,omitempty"`
	// Conditions is a collection of potential issues
	// +kubebuilder:validation:Optional
	// +listType=set
	Conditions []corev1.Condition `json:"conditions,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AccountManagement is the Schema for the accounts API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=accountmanagement,scope=Namespaced
type AccountManagement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountManagementSpec   `json:"spec,omitempty"`
	Status AccountManagementStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AccountManagementList contains a list of Account
type AccountManagementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AccountManagement `json:"items"`
}
