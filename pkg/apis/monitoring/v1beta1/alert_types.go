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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Alert contains the definition of a alert
type Alert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AlertSpec   `json:"spec,omitempty"`
	Status            AlertStatus `json:"status,omitempty"`
}

// AlertSpec specifies the details of a alert
// +k8s:openapi-gen=true
type AlertSpec struct {
	// Event is the raw event payload
	// +kubebuilder:validation:Optional
	Event string `json:"event,omitempty"`
	// Summary is human readable summary for the alert
	// +kubebuilder:validation:Required
	Summary string `json:"summary"`
}

// AlertStatus is the status of the alert
// +k8s:openapi-gen=true
type AlertStatus struct {
	// ArchivedAt is indicates if the alert has been archived
	// +kubebuilder:validation:Required
	ArchivedAt metav1.Time `json:"archivedAt"`
	// Detail provides a human readable message related to the current
	// status of the alert
	// +kubebuilder:validation:Optional
	Detail string `json:"detail,omitempty"`
	// Expiration is the time the silience will finish
	// +kubebuilder:validation:Optional
	Expiration metav1.Time `json:"expiration,omitempty"`
	// Status is the status of the alert
	// +kubebuilder:validation:Optional
	Status string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlertList contains a list of rules
type AlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Alert `json:"items"`
}
