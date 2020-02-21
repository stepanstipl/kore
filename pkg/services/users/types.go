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

import (
	"context"

	"github.com/appvia/kore/pkg/services/users/model"
)

const (
	// AuditCreate is a creation event
	AuditCreate = "create"
	// AuditDelete is a deletion event
	AuditDelete = "delete"
	// AuditUpdate is a update event
	AuditUpdate = "update"
)

// Config is the configurations for the store
type Config struct {
	// Driver is the database driver to use
	Driver string `json:"driver,omitempty"`
	// EnableLogging enables sql logging
	EnableLogging bool `json:"enable-logging,omitempty"`
	// StoreURL is the store endpoint url
	StoreURL string `json:"store-url,omitempty"`
}

// Interface defines the interface to the db store
type Interface interface {
	// Audit returns the audit interface
	Audit() Audit
	// Identities returns the identities interface
	Identities() Identities
	// Invitations returns the invitations interface
	Invitations() Invitations
	// Members returns the members interface
	Members() Members
	// Stop is called to shutdown the store and clean up resources
	Stop() error
	// Teams returns the teams interface
	Teams() Teams
	// Users returns the users interface
	Users() Users
}

// Audit is the interface to the audit service
type Audit interface {
	// Find is used to retrieve records from the log
	Find(context.Context, ...ListFunc) Find
	// Record records an event in the audit log
	Record(context.Context, ...AuditFunc) Log
	// Stop stops the service and releases resources
	Stop() error
}

type Find interface {
	// Do performs the query and returns the results
	Do() ([]*model.AuditEvent, error)
}

// RecordEntry defines a interface for adding a entry
type Log interface {
	// Event is responsible for entering the record into the audit
	Event(string)
}
