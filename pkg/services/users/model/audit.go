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

package model

import (
	"time"
)

// AuditEvent defines an audit event in the kore
type AuditEvent struct {
	// ID is the record id
	ID int `gorm:"primary_key"`
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Resource is the area of the API accessed in this audit operation (e.g. teams, ).
	Resource string
	// ResourceURI is the identifier of the resource in question.
	ResourceURI string
	// APIVersion is the version of the API in use for this operation.
	APIVersion string
	// Verb is the type of action performed (e.g. PUT, GET, etc)
	Verb string `gorm:"not null"`
	// Operation is the operation performed (e.g. UpdateCluster, CreateCluster, etc).
	Operation string
	// Team is the team whom event may be associated to
	Team string
	// User is the user which the event is related
	User string
	// StartedAt is the timestamp the operation was initiated
	StartedAt time.Time
	// CompletedAt is the timestamp the operation completed
	CompletedAt time.Time
	// ResponseCode indicates the HTTP status code of the operation (e.g. 200, 404, etc).
	ResponseCode int
	// Message is event message itself
	Message string `gorm:"not null"`
}
