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

	policy "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedPodSecurityPolicySpec defines the desired state of Cluster role
// +k8s:openapi-gen=true
type ManagedPodSecurityPolicySpec struct {
	// Clusters is used to apply the cluster role to a specific cluster
	// +kubebuilder:validation:Optional
	Clusters []corev1.Ownership `json:"clusters,omitempty"`
	// Teams is a filter on the teams
	// +kubebuilder:validation:Optional
	Teams []string `json:"teams,omitempty"`
	// Description describes the nature of this pod security policy
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Description string `json:"description,omitempty"`
	// Policy defined a managed pod security policy across the clusters
	// +kubebuilder:validation:Required
	Policy policy.PodSecurityPolicySpec `json:"policy,omitempty"`
}

// ManagedPodSecurityPolicyStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type ManagedPodSecurityPolicyStatus struct {
	// Conditions is a set of condition which has caused an error
	Conditions []corev1.Condition `json:"conditions"`
	// Status is overall status of the workspace
	Status corev1.Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubernetes is the Schema for the roles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=managedpodsecuritypoliies
type ManagedPodSecurityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedPodSecurityPolicySpec   `json:"spec,omitempty"`
	Status ManagedPodSecurityPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubernetesList contains a list of Managed
type ManagedPodSecurityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedPodSecurityPolicy `json:"items"`
}
