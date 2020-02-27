/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AuditEventSpec defines the desired state of User
// +k8s:openapi-gen=false
type AuditEventSpec struct {
	// CreatedAt is the timestamp of record creation
	CreatedAt metav1.Time `json:"createdAt,omitempty"`
	// Type is the type of event
	Type string `json:"type,omitempty"`
	// Team is the team whom event may be associated to
	Team string `json:"team,omitempty"`
	// User is the user which the event is related
	User string `json:"user,omitempty"`
	// Message is event message itself
	Message string `json:"message,omitempty"`
	// Resource is the name of the resource in question namespace/name
	Resource string `json:"resource,omitempty"`
	// ResourceUID is a unique id for the resource
	ResourceUID string `json:"resourceUID,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuditEvent is the Schema for the audit API
// +k8s:openapi-gen=false
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=auditevents
type AuditEvent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AuditEventSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuditEventList contains a list of audit event
type AuditEventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuditEvent `json:"items"`
}
