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

	// include the database drivers
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// storeImpl is the implementation of the users interface
type storeImpl struct {
	// dbc is the connection to the database
	dbc *gorm.DB
	// config is the configuration
	config Config
}

// New returns a db store to the consumer
func New(config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	log.WithField(
		"driver", config.Driver,
	).Info("initializing the database")

	db, err := gorm.Open(config.Driver, config.StoreURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(config.EnableLogging)
	db.DB().SetConnMaxLifetime(30 * time.Second)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(0)

	// @step: perform migrations on the models
	log.Info("performing database migrations")
	if err := model.Migrations(db); err != nil {
		return nil, err
	}

	return &storeImpl{dbc: db, config: config}, nil
}

// Audit returns the audit interface
func (s *storeImpl) Audit() Audit {
	return s
}

// Members returns the team members
func (s *storeImpl) Members() Members {
	return &membersImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Invitations returns the invitations interface
func (s *storeImpl) Invitations() Invitations {
	return &ivImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Teams returns the teams interface
func (s *storeImpl) Teams() Teams {
	return &teamImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Identities returns the identities interface
func (s *storeImpl) Identities() Identities {
	return &idImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Users returns the users interface
func (s *storeImpl) Users() Users {
	return &userImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Security returns the security interface
func (s *storeImpl) Security() Security {
	return &securityImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// Stop is called to free up resources
func (s *storeImpl) Stop() error {
	log.Info("shutting down the db store")
	if s.dbc != nil {
		return s.dbc.Close()
	}

	return nil
}

// Configs returns the configs interface
func (s *storeImpl) Configs() Configs {
	return &configImpl{
		Interface: s,
		conn:      s.dbc,
	}
}
