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
	"time"

	"github.com/appvia/kore/pkg/services/users/model"

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
	log.Info("perform database migrations")
	if err := model.Migrations(db); err != nil {
		return nil, err
	}

	return &storeImpl{dbc: db, config: config}, nil
}

// Audit retuns the audit interface
func (a *storeImpl) Audit() Audit {
	return a
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

// Stop is called to free up resources
func (s *storeImpl) Stop() error {
	log.Info("shutting down the db store")
	if s.dbc != nil {
		return s.dbc.Close()
	}

	return nil
}
