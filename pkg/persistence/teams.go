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

// Teams provides access to teams
type Teams interface {
	// AddUser is responsible for adding a user to the team
	AddUser(context.Context, string, string, []string) error
	// Delete removes a team from the store
	Delete(context.Context, *model.Team) error
	// Exists check if the team exists
	Exists(context.Context, string) (bool, error)
	// Get returns a team from the store
	Get(context.Context, string) (*model.Team, error)
	// List returns a list of teams from the store
	List(context.Context, ...ListFunc) ([]*model.Team, error)
	// Update updates a team in the store
	Update(context.Context, *model.Team) error
}

type teamImpl struct {
	Interface
	// load is for preload
	load []string
	// conn is the db connection
	conn *gorm.DB
}

// AddUser is responsible for adding a user to the team
func (t teamImpl) AddUser(ctx context.Context, username, team string, roles []string) error {
	u, err := t.Users().Get(ctx, username)
	if err != nil {
		return err
	}
	tm, err := t.Get(ctx, team)
	if err != nil {
		return err
	}

	return t.Members().Update(ctx, &model.Member{
		UserID: u.ID,
		TeamID: tm.ID,
		Roles:  roles,
	})
}

// Delete removes a team from the store
func (t teamImpl) Delete(ctx context.Context, team *model.Team) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return t.conn.Delete(team).Error
}

// Exists check if the team exists
func (t teamImpl) Exists(ctx context.Context, name string) (bool, error) {
	if _, err := t.Get(ctx, name); err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// Get returns a team from the store
func (t teamImpl) Get(ctx context.Context, name string) (*model.Team, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	team := &model.Team{}

	return team, t.conn.Where("name = ?", name).Find(team).Error
}

// List returns a list of teams from the store
func (t teamImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Team, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	filter := ApplyListOptions(opts...)

	q := Preload(t.load, t.conn).
		Model(&model.Team{}).
		Select("t.*").
		Table("teams t")

	if filter.HasTeam() {
		q = q.Where("t.name = ?", filter.GetTeam())
	}
	if filter.HasTeamID() {
		q = q.Where("t.id = ?", filter.GetTeamID())
	}

	var list []*model.Team

	return list, q.Find(&list).Error
}

// Update updates a team in the store
func (t teamImpl) Update(ctx context.Context, team *model.Team) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	if team.Name == "" {
		return errors.New("no team name defined")
	}

	return t.conn.
		Where("name = ?", team.Name).
		Assign(&model.Team{
			Name:        team.Name,
			Description: team.Description,
			Summary:     team.Summary,
			Identifier:  team.Identifier,
		}).
		FirstOrCreate(team).
		Error
}

// Preload allows access to the preload
func (t *teamImpl) Preload(v ...string) Teams {
	t.load = append(t.load, v...)

	return t
}
