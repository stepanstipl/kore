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

package audit

import "context"

import "github.com/appvia/kore/pkg/services/audit/model"

const (
	// Create is a creation event
	Create = "create"
	// Delete is a deletion event
	Delete = "delete"
	// Update is a update event
	Update = "update"
)

// Interface is the interface to the audit service
type Interface interface {
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
	Event(string) error
}

// Entry is the interface for the handler
type Entry interface {
	// Event records an event in the log
	Event(string) error
}

// Config provides the configuration for the audit service
type Config struct {
	// Driver is the database driver to use
	Driver string `json:"driver,omitempty"`
	// EnableLogging enables sql logging
	EnableLogging bool `json:"enable-logging,omitempty"`
	// StoreURL is the store endpoint url
	StoreURL string `json:"store-url,omitempty"`
}
