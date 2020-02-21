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
	"context"
	"errors"

	"github.com/appvia/kore/pkg/services/users/model"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// Members is the team members interface
type Members interface {
	// AddUser is responsible for adding a user to a team
	AddUser(context.Context, string, string, []string) error
	// Delete is responsible for removing a member from the team
	Delete(context.Context, *model.Member) error
	// DeleteBy removes based on a filter
	DeleteBy(context.Context, ...ListFunc) error
	// ListMembers returns a list of members in a team
	List(context.Context, ...ListFunc) ([]*model.Member, error)
	// Preload adds to the query preload
	Preload(...string) Members
	// Add is responsible for adding a member to a team
	Update(context.Context, *model.Member) error
}

// membersImpl implements the above interface
type membersImpl struct {
	Interface
	// load is for preloading
	load []string
	// conn is the db connection
	conn *gorm.DB
}

// AddUser is responsible for adding a user to a team
func (m *membersImpl) AddUser(ctx context.Context, user, team string, roles []string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	u, err := m.Users().Get(ctx, user)
	if err != nil {
		return err
	}
	t, err := m.Teams().Get(ctx, team)
	if err != nil {
		return err
	}

	return m.Members().Update(ctx, &model.Member{
		UserID: u.ID,
		TeamID: t.ID,
		Roles:  roles,
	})
}

// List returns a list of teams for a specific user
func (m *membersImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Member, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	var list []*model.Member

	q := Preload(m.load, m.conn).
		Select("m.*").
		Table("members m").
		Joins("LEFT JOIN teams t ON t.id = m.team_id").
		Joins("LEFT JOIN users u ON u.id = m.user_id")

	if terms.HasTeam() {
		q = q.Where("t.name = ?", terms.GetTeam())
	}
	if terms.HasUser() {
		q = q.Where("u.username = ?", terms.GetUser())
	}
	if err := q.Preload("Team").Find(&list).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}

		return []*model.Member{}, nil
	}

	return list, nil
}

// DeleteBy is responsible for deleting by filter
func (m *membersImpl) DeleteBy(ctx context.Context, filters ...ListFunc) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	if len(filters) <= 0 {
		return errors.New("no filters for delete by on users")
	}

	terms := ApplyListOptions(filters...)

	q := m.conn.
		Model(&model.Member{}).
		Select("m.*").
		Table("members m").
		Joins("JOIN teams t ON t.id = m.team_id").
		Joins("JOIN users u ON u.id = m.user_id")

	if terms.HasUser() {
		q = q.Where("u.username = ?", terms.GetUser())
	}
	if terms.HasTeam() {
		q = q.Where("t.name = ?", terms.GetTeam())
	}

	list := []*model.Member{}
	if err := q.Find(&list).Error; err != nil {
		return err
	}

	for _, x := range list {
		if err := m.conn.Model(&model.Member{}).Delete(x).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete is responsible for removing a member from the team
func (m *membersImpl) Delete(ctx context.Context, member *model.Member) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	if member.UserID == 0 {
		return errors.New("no user id defined")
	}
	if member.TeamID == 0 {
		return errors.New("no team id defined")
	}

	return m.conn.
		Where("user_id = ?", member.UserID).
		Where("team_id = ?", member.TeamID).
		Delete(member).
		Error
}

// Add is responsible for adding a member to a team
func (m membersImpl) Update(ctx context.Context, member *model.Member) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return m.conn.FirstOrCreate(member, member).Error
}

// Preload adds proloading to the queries
func (m *membersImpl) Preload(v ...string) Members {
	m.load = append(m.load, v...)

	return m
}
