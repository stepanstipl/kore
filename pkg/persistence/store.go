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
	"database/sql"
	"time"

	"github.com/appvia/kore/pkg/persistence/migrations"
	"github.com/appvia/kore/pkg/persistence/model"

	// include the database drivers
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"

	bindata "github.com/golang-migrate/migrate/source/go_bindata"
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
	// @TODO once we are happy with the go-migrate we should turn
	// gorm automigate off and use that only
	log.Info("performing gorm database migrations")
	if err := model.Migrations(db); err != nil {
		return nil, err
	}

	log.Info("performing go-migrate database migrations")
	if err := Migrations("kore", config.Driver, config.StoreURL); err != nil {
		if err != migrate.ErrNoChange {
			return nil, err
		}
	}

	return &storeImpl{dbc: db, config: config}, nil
}

// Migrations is responsible for performing the database migrations
func Migrations(dbname, driver, databaseURL string) error {
	files := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	// @step: perform the go-migrate migrations
	source, err := bindata.WithInstance(files)
	if err != nil {
		return err
	}

	db, err := sql.Open(driver, databaseURL)
	if err != nil {
		return err
	}

	dest, err := mysql.WithInstance(db, &mysql.Config{
		DatabaseName:    dbname,
		MigrationsTable: "migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("kore", source, "kore", dest)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		log.WithError(err).Error("trying to perform database migrations")

		return err
	}

	return nil
}

// Alerts returns the alerts interface
func (s *storeImpl) Alerts() Alerts {
	return &alertsImpl{
		Interface: s,
		conn:      s.dbc,
	}
}

// AlertRules returns the alerts interface
func (s *storeImpl) AlertRules() AlertRules {
	return &arulesImpl{
		Interface: s,
		conn:      s.dbc,
	}
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

func (s storeImpl) IsNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
