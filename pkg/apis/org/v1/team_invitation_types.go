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

// TeamInvitationSpec defines the desired state of Team
// +k8s:openapi-gen=true
type TeamInvitationSpec struct {
	// Username is the user being bound to the team
	// +kubebuilder:validation:Required
	Username string `json:"username"`
	// Team is the name of the team
	// +kubebuilder:validation:Required
	Team string `json:"team"`
}

// TeamInvitationStatus defines the observed state of Team
// +k8s:openapi-gen=true
type TeamInvitationStatus struct {
	// Conditions is a collection of possible errors
	// +kubebuilder:validation:Optional
	// +listType=set
	Conditions []corev1.Condition `json:"conditions"`
	// Status is the status of the resource
	Status corev1.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeamInvitation is the Schema for the teams API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=teaminvitations
type TeamInvitation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamInvitationSpec   `json:"spec,omitempty"`
	Status TeamInvitationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TeamInvitationList contains a list of Team
type TeamInvitationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TeamInvitation `json:"items"`
}
