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

package audit

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/services/audit/model"

	// include the database drivers
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type auditImpl struct {
	// dbc is the connection to the database
	dbc *gorm.DB
	// config is the configuration
	config Config
}

// New creates and returns an audit service
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
	db.DB().SetMaxIdleConns(4)
	db.DB().SetMaxOpenConns(0)

	// @step: perform migrations on the models
	log.Info("perform database migrations")
	if err := model.Migrations(db); err != nil {
		return nil, err
	}

	return &auditImpl{dbc: db, config: config}, nil
}

// Find is used to return results from the log
func (a *auditImpl) Find(ctx context.Context, filters ...ListFunc) Find {
	return newQuery(ctx, a.dbc, filters...)
}

// Record is responsible for adding an entry into the log
func (a *auditImpl) Record(ctx context.Context, fields ...AuditFunc) Log {
	return newEntry(ctx, a.dbc, fields...)
}

// Stop is used to free up resources
func (a *auditImpl) Stop() error {
	return nil
}
