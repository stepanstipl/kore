/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallPlanSpec defines the desired state of Allocation
// +k8s:openapi-gen=true
type InstallPlanSpec struct {
	// Approved indicates if the update has been approved
	Approved bool `json:"approved,omitempty"`
}

// InstallPlanStatus defines the observed state of Allocation
// +k8s:openapi-gen=true
type InstallPlanStatus struct {
	// Conditions is a collection of potential issues
	// +listType
	Conditions []corev1.Condition `json:"conditions,omitempty"`
	// Deployed is the applciation deployment parameters
	Deployed AppDeployment `json:"deployed"`
	// Update is the incoming deployment is requiring approval
	Update AppDeployment `json:"update,omitempty"`
	// Status is the general status of the resource
	Status corev1.Status `json:"status,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstallPlan is the Schema for the allocations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=installplans,scope=Namespaced
type InstallPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstallPlanSpec   `json:"spec,omitempty"`
	Status InstallPlanStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstallPlanList contains a list of Allocation
type InstallPlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstallPlan `json:"items"`
}
