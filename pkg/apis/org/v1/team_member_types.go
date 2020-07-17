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

// TeamMemberSpec defines the desired state of Team
// +k8s:openapi-gen=true
type TeamMemberSpec struct {
	// Role is the role of the user in the team
	// +kubebuilder:validation:Required
	Roles []string `json:"roles"`
	// Team is the name of the team
	// +kubebuilder:validation:Required
	Team string `json:"team"`
	// Username is the user being bound to the team
	// +kubebuilder:validation:Required
	Username string `json:"username"`
}

// TeamMemberStatus defines the observed state of Team
// +k8s:openapi-gen=true
type TeamMemberStatus struct {
	// Conditions is a collection of possible errors
	// +kubebuilder:validation:Optional
	Conditions []corev1.Condition `json:"conditions"`
	// Status is the status of the resource
	Status corev1.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeamMember is the Schema for the teams API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=members
type TeamMember struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamMemberSpec   `json:"spec,omitempty"`
	Status TeamMemberStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeamMemberList contains a list of Team
type TeamMemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TeamMember `json:"items"`
}
