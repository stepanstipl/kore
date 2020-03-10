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

package users

import "github.com/appvia/kore/pkg/services/users/model"

// AuditFunc sets an option in the record
type AuditFunc func(m *model.AuditEvent)

// Type sets the type of the event
func Type(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Type = v
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

// Resource sets the resource
func Resource(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.Resource = v
	}
}

// Resource sets the resource uid
func ResourceUID(v string) AuditFunc {
	return func(m *model.AuditEvent) {
		m.ResourceUID = v
	}
}
