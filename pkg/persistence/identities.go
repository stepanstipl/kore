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
	"strings"

	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// Identities provides access to the identities
type Identities interface {
	// Delete removes a user identity from the kore
	Delete(context.Context, *model.Identity) error
	// Get returns a single identity if any
	Get(context.Context, ...ListFunc) (*model.Identity, error)
	// List returns a filtered list of user identities
	List(context.Context, ...ListFunc) ([]*model.Identity, error)
	// Preload allows the user to preload
	Preload(...string) Identities
	// Update adds or updates an user identity
	Update(context.Context, *model.Identity) error
}

type idImpl struct {
	Interface
	// load provides a preloading cap
	load []string
	// conn is the database connection
	conn *gorm.DB
}

// Delete removes a user identity from the kore
func (i *idImpl) Delete(ctx context.Context, id *model.Identity) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	if id.UserID == 0 {
		return errors.New("identity has no user id")
	}
	if id.Provider == "" {
		return errors.New("no provider name defined")
	}

	return Preload(i.load, i.conn).
		Where("user_id = ? AND provider = ?", id.UserID, id.Provider).
		Delete(id).
		Error
}

// Get returns a single identity if any
func (i *idImpl) Get(ctx context.Context, opts ...ListFunc) (*model.Identity, error) {
	list, err := i.List(ctx, opts...)
	if err != nil {
		return nil, err
	}
	switch len(list) {
	case 0:
		return nil, gorm.ErrRecordNotFound
	case 1:
		return list[0], nil
	default:
		return nil, errors.New("filter matched more than one record")
	}
}

// List returns a filtered list of user identities
func (i *idImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Identity, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	q := Preload(i.load, i.conn).
		Preload("User").
		Select("i.*").
		Table("identities i").
		Joins("JOIN users u ON u.id = i.user_id")

	if terms.HasUser() {
		q = q.Where("u.username = ?", terms.GetUser())
	}
	if terms.HasProviders() {
		q = q.Where("i.provider IN ?", terms.GetProviders())
	}
	if terms.HasProvider() {
		q = q.Where("i.provider = ?", terms.GetProvider())
	}

	list := []*model.Identity{}

	if err := q.Find(&list).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}

		return list, nil
	}

	return list, nil
}

// Update adds or updates an user identity
func (i *idImpl) Update(ctx context.Context, id *model.Identity) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	if id.UserID == 0 {
		return errors.New("user id is not defined")
	}
	if id.Provider == "" {
		return errors.New("provider is not defined")
	}
	if id.ProviderToken != "" && !strings.HasPrefix("md5", id.ProviderToken) {
		//encoded, err := bcrypt.GenerateFromPassword([]byte(id.ProviderToken), 9)
		// 38ms vs 1us
		id.ProviderToken = "md5:" + utils.HashString(id.ProviderToken)
	}

	return i.conn.
		Where("user_id = ? AND provider = ?", id.UserID, id.Provider).
		Assign(&model.Identity{
			Extras:           id.Extras,
			Provider:         id.Provider,
			ProviderUsername: id.ProviderUsername,
			ProviderEmail:    id.ProviderEmail,
			ProviderToken:    id.ProviderToken,
			ProviderUID:      id.ProviderUID,
			UserID:           id.UserID,
		}).
		FirstOrCreate(id).
		Error
}

// Preload allows the user to preload
func (i *idImpl) Preload(p ...string) Identities {
	i.load = append(i.load, p...)

	return i
}
