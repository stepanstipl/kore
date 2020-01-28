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

// Invitations provides access to the invitations
type Invitations interface {
	// Get retrieves a invitation by filter
	Get(context.Context, ...ListFunc) (*model.Invitation, error)
	// Delete removes an invitation
	Delete(context.Context, *model.Invitation) error
	// DeleteBy removes an invitation by filter
	DeleteBy(context.Context, ...ListFunc) error
	// List returns a filtered list of invitations
	List(context.Context, ...ListFunc) ([]*model.Invitation, error)
	// Preload allows for the consumer to select the preloaded fields
	Preload(...string) Invitations
	// Update updates or creates and invitations
	Update(context.Context, *model.Invitation) error
}

type ivImpl struct {
	Interface
	// load is the preloaded fields
	load []string
	// conn is the db connection for this query
	conn *gorm.DB
}

// Delete removes an invitation
func (i *ivImpl) Delete(ctx context.Context, iv *model.Invitation) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return i.conn.Delete(iv).Error
}

// DeleteBy removes invitatons by filter
func (i *ivImpl) DeleteBy(ctx context.Context, filters ...ListFunc) error {
	if len(filters) <= 0 {
		return errors.New("no filters defined for deletion of invitations")
	}

	terms := ApplyListOptions(filters...)
	query := i.conn.
		Model(&model.Invitation{}).
		Select("i.*").
		Table("invitations i").
		Joins("JOIN teams t ON t.id = i.team_id").
		Joins("JOIN users u ON u.id = i.user_id")

	if terms.HasUser() {
		query = query.Where("u.username = ?", terms.GetUser())
	}
	if terms.HasTeam() {
		query = query.Where("t.name = ?", terms.GetTeam())
	}

	// @TODO needs fixing up
	list := []*model.Invitation{}
	err := query.Find(&list).Error
	if err != nil {
		return err
	}

	for _, x := range list {
		if err := i.conn.Model(&model.Invitation{}).Delete(x).Error; err != nil {
			return err
		}
	}

	return nil
}

// Get returns a invitation by filter - else a no record found error
func (i *ivImpl) Get(ctx context.Context, opts ...ListFunc) (*model.Invitation, error) {
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
		return nil, errors.New("matched more than one record")
	}
}

// List returns a filtered list of invitations
func (i *ivImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Invitation, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	q := Preload(i.load, i.conn).
		Select("i.*").
		Table("invitations i").
		Joins("JOIN teams t ON t.id = i.team_id").
		Joins("JOIN users u ON u.id = i.user_id")

	if terms.HasTeam() {
		q = q.Where("t.name = ?", terms.GetTeam())
	}
	if terms.HasUser() {
		q = q.Where("u.username = ?", terms.GetUser())
	}

	var list []*model.Invitation

	if err := q.Find(&list).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}

		return []*model.Invitation{}, nil
	}

	return list, nil
}

// Update updates or creates and invitations
func (i *ivImpl) Update(ctx context.Context, iv *model.Invitation) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	if iv.Team != nil {
		iv.TeamID = iv.Team.ID
	}
	if iv.User != nil {
		iv.UserID = iv.User.ID
	}

	if iv.TeamID == 0 {
		return errors.New("no team id defined")
	}
	if iv.UserID == 0 {
		return errors.New("no user id defined")
	}

	return i.conn.
		Set("gorm:save_associations", false).
		Assign(&model.Invitation{
			TeamID: iv.TeamID,
			UserID: iv.UserID,
		}).
		FirstOrCreate(iv, &model.Invitation{TeamID: iv.ID, UserID: iv.UserID}).
		Error
}

// Preload allows for the consumer to select the preloaded fields
func (i *ivImpl) Preload(v ...string) Invitations {
	i.load = append(i.load, v...)

	return i
}
