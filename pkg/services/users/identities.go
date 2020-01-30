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

package users

import (
	"context"
	"errors"

	"github.com/appvia/kore/pkg/services/users/model"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// Identities provides access to the identities
type Identities interface {
	// Delete removes a user identity from the hub
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

// Delete removes a user identity from the hub
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

// List returns a filtered list of user identites
func (i *idImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Identity, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	q := Preload(i.load, i.conn).
		Select("i.*").
		Table("identities i").
		Joins("JOIN users u ON u.id = i.user_id")

	if terms.HasUser() {
		q = q.Where("u.username = ?", terms.GetUser())
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
