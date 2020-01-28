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

package model

import "time"

// AuditEvent defines an audit event in the hub
type AuditEvent struct {
	// ID is the record id
	ID int `gorm:"primary_key"`
	// CreatedAt is the timestamp of record creation
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp"`
	// Type is the type of event
	Type string `gorm:"not null"`
	// Team is the team whom event may be associated to
	Team string
	// User is the user which the event is related
	User string
	// Message is event message itself
	Message string `gorm:"not null"`
	// Resource is the name of the resource in question namespace/name
	Resource string
	// ResourceUID is a unique id for the resource
	ResourceUID string
}
