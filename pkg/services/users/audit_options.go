/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
