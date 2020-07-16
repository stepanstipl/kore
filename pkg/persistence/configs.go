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
	"errors"

	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// Configs defines the config interface to the store
type Configs interface {
	// Delete removes a config from the store
	Delete(context.Context, *model.Config) (*model.Config, error)
	// Exists check if config exists
	Exists(context.Context, string) (bool, error)
	// Get returns a config from the store
	Get(context.Context, string) (*model.Config, error)
	// List returns a list of all configs from the store
	List(context.Context, ...ListFunc) ([]*model.Config, error)
	// // Update updates a config in the store
	Update(context.Context, *model.Config) error
}

// configImpl handles access to the config model
type configImpl struct {
	Interface

	conn *gorm.DB
}

// Delete removes a config from the store
func (v configImpl) Delete(ctx context.Context, config *model.Config) (*model.Config, error) {

	if config.Name == "" {
		return nil, errors.New("invalid name for deletion")
	}

	q := v.conn

	return config, q.Delete(config).Error
}

// Exists check if config exists
func (v configImpl) Exists(ctx context.Context, name string) (bool, error) {
	if _, err := v.Get(ctx, name); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Get returns a config from the store
func (v configImpl) Get(ctx context.Context, name string) (*model.Config, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	record := &model.Config{}

	q := v.conn.Preload("Items")

	err := q.Where("name = ?", name).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return record, nil
}

// List returns a list of all configs from the store
func (v configImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Config, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	q := v.conn
	var list []*model.Config

	err := q.Preload("Items").Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// // Update updates a config in the store
func (v configImpl) Update(ctx context.Context, config *model.Config) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return v.conn.Model(&model.Config{}).
		Where("name = ?", config.Name).
		Assign(&model.Config{
			Name:  config.Name,
			Items: config.Items,
		}).
		FirstOrCreate(config).
		Error
}
