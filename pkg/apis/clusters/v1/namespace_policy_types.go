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

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamepacePolicySpec defines the desired state of Cluster role
// +k8s:openapi-gen=true
type NamepacePolicySpec struct {
	// DefaultLabels are the labels applied to all managed namespaces
	// +kubebuilder:validation:Optional
	DefaultLabels map[string]string `json:"defaultLabels,omitempty"`
	// DefaultAnnotations are default annotations applied to all managed namespaces
	// +kubebuilder:validation:Optional
	DefaultAnnotations map[string]string `json:"defaultAnnotations,omitempty"`
	// DefaultLimits are the default resource limits applied to the namespace
	// +kubebuilder:validation:Optional
	DefaultLimits *core.LimitRange `json:"defaultLimits,omitempty"`
}

// NamepacePolicyStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type NamepacePolicyStatus struct {
	// Conditions is a set of condition which has caused an error
	// +listType=set
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubernetes is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=namespacepolicy
type NamepacePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamepacePolicySpec   `json:"spec,omitempty"`
	Status NamepacePolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList contains a list of Managed
type NamepacePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamepacePolicy `json:"items"`
}
