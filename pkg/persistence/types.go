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
	"context"

	"github.com/appvia/kore/pkg/persistence/model"
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
	// Alerts returns the alerts interface
	Alerts() Alerts
	// AlertRules returns the alert rules interface
	AlertRules() AlertRules
	// Audit returns the audit interface
	Audit() Audit
	// Identities returns the identities interface
	Identities() Identities
	// Invitations returns the invitations interface
	Invitations() Invitations
	// Members returns the members interface
	Members() Members
	// Security returns the security interface
	Security() Security
	// Stop is called to shutdown the store and clean up resources
	Stop() error
	// Teams returns the teams interface
	Teams() Teams
	// TeamAssets returns the team assets interface
	TeamAssets() TeamAssets
	// Users returns the users interface
	Users() Users
	// Config returns the config interface
	Configs() Configs
	// IsNotFound checks if the supplied error is a not found error
	IsNotFound(err error) bool
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

// Find is an action interface
type Find interface {
	// Do performs the query and returns the results
	Do() ([]*model.AuditEvent, error)
}

// Log defines a interface for adding a entry
type Log interface {
	// Event is responsible for entering the record into the audit
	Event(string)
}
