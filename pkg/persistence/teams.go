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
	"time"

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
	// RecordTeamIdentity persists a new team identifier
	RecordTeamIdentity(ctx context.Context, teamIdentifier string, teamName string) error
	// MarkTeamIdentityDeleted marks a specific team identifier as deleted
	MarkTeamIdentityDeleted(ctx context.Context, teamIdentifier string) error
	// RecordAsset records an asset as being owned by a team
	RecordAsset(ctx context.Context, teamIdentifier string, assetIdentifier string, assetType model.TeamAssetType, assetName string) error
	// GetAsset retrieves details of an asset from the store
	GetAsset(ctx context.Context, teamIdentifier string, assetIdentifier string) (*model.TeamAsset, error)
	// MarkAssetDeleted records an asset as no longer being active
	MarkAssetDeleted(ctx context.Context, teamIdentifier string, assetIdentifier string) error
	// MarkAssetUndeleted records an asset as being active after previously being deleted
	MarkAssetUndeleted(ctx context.Context, teamIdentifier string, assetIdentifier string, assetName string) error
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

// RecordTeamIdentiy records the existence of a team
func (t teamImpl) RecordTeamIdentity(ctx context.Context, teamIdentifier string, teamName string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Create(&model.TeamIdentity{
		TeamIdentifier: teamIdentifier,
		TeamName:       teamName,
	}).Error
}

func (t teamImpl) MarkTeamIdentityDeleted(ctx context.Context, teamIdentifier string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Exec(
		"UPDATE team_identities SET deleted_at = ? WHERE team_identifier = ?",
		time.Now(),
		teamIdentifier,
	).Error
}

// RecordAsset records an asset as being owned by a team
func (t teamImpl) RecordAsset(ctx context.Context, teamIdentifier string, assetIdentifier string, assetType model.TeamAssetType, assetName string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Create(&model.TeamAsset{
		TeamIdentifier:  teamIdentifier,
		AssetIdentifier: assetIdentifier,
		AssetType:       assetType,
		AssetName:       assetName,
	}).Error
}

// MarkAssetDeleted records an asset as no longer being active
func (t teamImpl) MarkAssetDeleted(ctx context.Context, teamIdentifier string, assetIdentifier string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Exec(
		"UPDATE team_assets SET deleted_at = ? WHERE team_identifier = ? and asset_identifier = ?",
		time.Now(),
		teamIdentifier,
		assetIdentifier,
	).Error
}

// MarkAssetUndeleted records an asset as being active after previously being deleted
func (t teamImpl) MarkAssetUndeleted(ctx context.Context, teamIdentifier string, assetIdentifier string, assetName string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Exec(
		"UPDATE team_assets SET deleted_at = NULL, asset_name = ? WHERE team_identifier = ? and asset_identifier = ?",
		assetName,
		teamIdentifier,
		assetIdentifier,
	).Error
}

func (t teamImpl) GetAsset(ctx context.Context, teamIdentifier string, assetIdentifier string) (*model.TeamAsset, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	asset := &model.TeamAsset{}

	// @note: uses unscoped to include rows with deleted_at set - gorm by default
	// excludes 'soft deleted' rows - http://gorm.io/docs/delete.html#Soft-Delete
	return asset, t.conn.Unscoped().Where("team_identifier = ? AND asset_identifier = ?", teamIdentifier, assetIdentifier).First(asset).Error
}

// Preload allows access to the preload
func (t *teamImpl) Preload(v ...string) Teams {
	t.load = append(t.load, v...)

	return t
}
