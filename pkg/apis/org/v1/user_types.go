/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
	Disabled bool `json:"disabled,omitempty"`
	// Email is the email for the user
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Email string `json:"email,omitempty"`
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
	// +listType
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
