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

// UserSpec defines the desired state of User
// +k8s:openapi-gen=true
type UserSpec struct {
	// Disabled indicates if the user is disabled
	Disabled bool `json:"disabled"`
	// Email is the email for the user
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Email string `json:"email"`
	// Username is the userame or identity for this user
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Username string `json:"username"`
}

// UserStatus defines the observed state of User
// +k8s:openapi-gen=true
type UserStatus struct {
	// Conditions is collection of potentials error causes
	// +kubebuilder:validation:Optional
	Conditions []corev1.Condition `json:"conditions"`
	// Status provides an overview of the user status
	Status corev1.Status `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// User is the Schema for the users API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=users
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec,omitempty"`
	Status UserStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UserList contains a list of User
type UserList struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}
