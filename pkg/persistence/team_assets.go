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
	"time"

	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// TeamAssets provides access to team assets
type TeamAssets interface {
	// RecordTeamIdentity persists a new team identifier
	RecordTeamIdentity(ctx context.Context, teamIdentifier string, teamName string) error
	// GetTeamNameForIdentity returns the team name for the supplied idenfitier
	GetTeamNameForIdentity(ctx context.Context, teamIdentifier string) (string, error)
	// MarkTeamIdentityDeleted marks a specific team identifier as deleted
	MarkTeamIdentityDeleted(ctx context.Context, teamIdentifier string) error
	// RecordAsset records an asset as being owned by a team
	RecordAsset(ctx context.Context, teamIdentifier string, assetIdentifier string, assetType model.TeamAssetType, assetName string, provider string) error
	// GetAsset retrieves details of an asset from the store
	GetAsset(ctx context.Context, teamIdentifier string, assetIdentifier string) (*model.TeamAsset, error)
	// MarkAssetDeleted records an asset as no longer being active
	MarkAssetDeleted(ctx context.Context, teamIdentifier string, assetIdentifier string) error
	// MarkAssetUndeleted records an asset as being active after previously being deleted
	MarkAssetUndeleted(ctx context.Context, teamIdentifier string, assetIdentifier string, assetName string, provider string) error
	// ListAssets returns a list of assets filtered by the supplied filters
	ListAssets(ctx context.Context, filters ...TeamAssetFilterFunc) ([]*model.TeamAsset, error)
	// StoreAssetCost persists a new asset cost record
	StoreAssetCost(ctx context.Context, cost *model.TeamAssetCost) error
	// ListCosts returns a list asset costs filtered bythe supplied filters
	ListCosts(ctx context.Context, filters ...TeamAssetFilterFunc) ([]*model.TeamAssetCost, error)
}

type teamAssetsImpl struct {
	Interface
	// conn is the db connection
	conn *gorm.DB
}

// RecordTeamIdentiy records the existence of a team
func (t teamAssetsImpl) RecordTeamIdentity(ctx context.Context, teamIdentifier string, teamName string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Create(&model.TeamIdentity{
		TeamIdentifier: teamIdentifier,
		TeamName:       teamName,
	}).Error
}

func (t teamAssetsImpl) GetTeamNameForIdentity(ctx context.Context, teamIdentifier string) (string, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	teamIdentity := &model.TeamIdentity{}

	return teamIdentity.TeamName, t.conn.Where("team_identifier = ?", teamIdentifier).Find(teamIdentity).Error
}

func (t teamAssetsImpl) MarkTeamIdentityDeleted(ctx context.Context, teamIdentifier string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Exec(
		"UPDATE team_identities SET deleted_at = ? WHERE team_identifier = ?",
		time.Now(),
		teamIdentifier,
	).Error
}

// RecordAsset records an asset as being owned by a team
func (t teamAssetsImpl) RecordAsset(ctx context.Context, teamIdentifier string, assetIdentifier string, assetType model.TeamAssetType, assetName string, provider string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Create(&model.TeamAsset{
		TeamIdentifier:  teamIdentifier,
		AssetIdentifier: assetIdentifier,
		AssetType:       assetType,
		AssetName:       assetName,
		Provider:        provider,
	}).Error
}

// MarkAssetDeleted records an asset as no longer being active
func (t teamAssetsImpl) MarkAssetDeleted(ctx context.Context, teamIdentifier string, assetIdentifier string) error {
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
func (t teamAssetsImpl) MarkAssetUndeleted(ctx context.Context, teamIdentifier string, assetIdentifier string, assetName string, provider string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Exec(
		"UPDATE team_assets SET deleted_at = NULL, asset_name = ?, provider = ? WHERE team_identifier = ? and asset_identifier = ?",
		assetName,
		provider,
		teamIdentifier,
		assetIdentifier,
	).Error
}

func (t teamAssetsImpl) GetAsset(ctx context.Context, teamIdentifier string, assetIdentifier string) (*model.TeamAsset, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	asset := &model.TeamAsset{}

	// @note: uses unscoped to include rows with deleted_at set - gorm by default
	// excludes 'soft deleted' rows - http://gorm.io/docs/delete.html#Soft-Delete
	return asset, t.conn.Unscoped().Where("team_identifier = ? AND asset_identifier = ?", teamIdentifier, assetIdentifier).First(asset).Error
}

func (t teamAssetsImpl) ListAssets(ctx context.Context, filters ...TeamAssetFilterFunc) ([]*model.TeamAsset, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	q := t.conn.Model(&model.TeamAsset{})
	q = t.applyAssetListFilters(q, filters)

	var list []*model.TeamAsset
	return list, q.Find(&list).Error
}

func (t teamAssetsImpl) StoreAssetCost(ctx context.Context, cost *model.TeamAssetCost) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return t.conn.Create(cost).Error
}

func (t teamAssetsImpl) ListCosts(ctx context.Context, filters ...TeamAssetFilterFunc) ([]*model.TeamAssetCost, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	q := t.conn.Preload("Team").Preload("Asset")
	q = t.applyCostListFilters(q, filters)

	var list []*model.TeamAssetCost
	return list, q.Find(&list).Error
}

func (t teamAssetsImpl) applyCostListFilters(q *gorm.DB, filters []TeamAssetFilterFunc) *gorm.DB {
	filter := ApplyTeamAssetListOptions(filters...)

	q = t.applySharedTeamAssetFilters(q, filter)

	if filter.From != nil {
		q = q.Where("usage_start_time >= ?", filter.From)
	}

	if filter.To != nil {
		q = q.Where("usage_start_time <= ?", filter.To)
	}

	if filter.BillingYear != nil && filter.BillingMonth != nil {
		q = q.Where("billing_year = ? and billing_month = ?", filter.BillingYear, filter.BillingMonth)
	}

	if filter.Account != "" {
		q = q.Where("account = ?", filter.Account)
	}

	return q
}

func (t teamAssetsImpl) applyAssetListFilters(q *gorm.DB, filters []TeamAssetFilterFunc) *gorm.DB {
	filter := ApplyTeamAssetListOptions(filters...)

	if filter.WithDeleted {
		// @note: uses unscoped to include rows with deleted_at set - gorm by default
		// excludes 'soft deleted' rows - http://gorm.io/docs/delete.html#Soft-Delete
		q = q.Unscoped()
	}

	q = t.applySharedTeamAssetFilters(q, filter)

	return q
}

func (t teamAssetsImpl) applySharedTeamAssetFilters(q *gorm.DB, filter *TeamAssetListOptions) *gorm.DB {
	if filter.TeamIdentifier != "" {
		q = q.Where("team_identifier = ?", filter.TeamIdentifier)
	}

	if filter.AssetIdentifier != "" {
		q = q.Where("asset_identifier = ?", filter.AssetIdentifier)
	}

	if filter.Provider != "" {
		q = q.Where("provider = ?", filter.Provider)
	}

	return q
}
