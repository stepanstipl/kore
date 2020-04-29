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

package persistence

import (
	"time"

	"github.com/appvia/kore/pkg/persistence/model"
)

// AuditFunc sets an option in the record
type AuditFunc func(m *model.AuditEvent)

// Resource sets the area of the API accessed in this audit operation (e.g. teams, plans, etc).
func Resource(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Resource = v
	}
}

// ResourceURI sets the URI of the resource being operated on.
func ResourceURI(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.ResourceURI = v
	}
}

// APIVersion sets the API version for the operation which caused the audit.
func APIVersion(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.APIVersion = v
	}
}

// Verb sets the verb (e.g. GET/POST/PUT etc) of the event
func Verb(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Verb = v
	}
}

// Operation sets the operation (e.g. CreateCluster, GetTeams etc) of the event
func Operation(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Operation = v
	}
}

// Team sets the team in the event
func Team(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Team = v
	}
}

// User sets the user in the event
func User(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.User = v
	}
}

// StartedAt sets the time the event started
func StartedAt(v time.Time) AuditFunc {
	return func(m *model.AuditEvent) {
		m.StartedAt = v
	}
}

// CompletedAt sets the time the event completed
func CompletedAt(v time.Time) AuditFunc {
	return func(m *model.AuditEvent) {
		m.CompletedAt = v
	}
}

// ResponseCode sets the resulting HTTP status code of this operation.
func ResponseCode(v int) AuditFunc {
	return func(m *model.AuditEvent) {
		m.ResponseCode = v
	}
}
