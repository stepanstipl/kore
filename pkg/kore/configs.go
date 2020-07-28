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

package kore

import (
	"context"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/utils/validation"
	log "github.com/sirupsen/logrus"
)

// Configs is the kore api configs interface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Configs
type Configs interface {
	// Delete removes a config from the kore
	Delete(context.Context, string) (*configv1.Config, error)
	// Exists checks if the config exists
	Exists(context.Context, string) (bool, error)
	// Get returns the config from the kore
	Get(context.Context, string) (*configv1.Config, error)
	// List returns a list of stored configs
	List(context.Context) (*configv1.ConfigList, error)
	// Update is responsible for updating the config
	Update(context.Context, *configv1.Config) (*configv1.Config, error)
}

type configImpl struct {
	*hubImpl
}

// Delete removes a config from the kore
func (v *configImpl) Delete(ctx context.Context, key string) (*configv1.Config, error) {
	u, err := v.persistenceMgr.Configs().Get(ctx, key)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Name": u.Name,
	}).Info("deleting the config from the kore")

	if _, err := v.persistenceMgr.Configs().Delete(ctx, u); err != nil {
		log.WithError(err).Error("trying to remove config key from kore")

		return nil, err
	}

	return DefaultConvertor.FromConfigModel(u), nil
}

// Exists checks if the config exists
func (v *configImpl) Exists(ctx context.Context, key string) (bool, error) {
	return v.persistenceMgr.Configs().Exists(ctx, key)
}

// Get returns the config from the kore
func (v *configImpl) Get(ctx context.Context, key string) (*configv1.Config, error) {
	configkey, err := v.persistenceMgr.Configs().Get(ctx, key)
	if err != nil {
		if persistence.IsNotFound(err) {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("trying to retrieve specific config key")

		return nil, err
	}

	return DefaultConvertor.FromConfigModel(configkey), nil
}

// List returns a list of stored configs
func (v *configImpl) List(ctx context.Context) (*configv1.ConfigList, error) {
	list, err := v.persistenceMgr.Configs().List(ctx)
	if err != nil {
		log.WithError(err).Error("trying to retrieve list")

		return nil, err
	}

	return DefaultConvertor.FromConfigModelList(list), err
}

// Update is responsible for updating the config
func (v *configImpl) Update(ctx context.Context, config *configv1.Config) (*configv1.Config, error) {
	valErr := validation.NewError("config has failed validation")
	if config.Name == "" {
		valErr.AddFieldError("Name", validation.Required, "can not be empty")
	}
	if config.Spec.Values == nil {
		valErr.AddFieldError("Values", validation.Required, "can not be empty")
	}
	if valErr.HasErrors() {
		return nil, valErr
	}

	if err := v.persistenceMgr.Configs().Update(ctx, DefaultConvertor.ToConfigModel(config)); err != nil {
		log.WithError(err).Error("trying to update the config in kore")

		return nil, err
	}

	return config, nil
}
