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
	"fmt"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils"
)

// TeamAssets represents the interface to managing team assets.
type TeamAssets interface {
	// EnsureTeamIdentifier validates and/or assigns a team identifier - if suppliedTeamIdentifier is non-empty,
	// it is checked against the team's current identifier and the return indicates whether the identifier is
	// valid. If the team has no identifier one will be assigned and returned.
	EnsureTeamIdentifier(ctx context.Context, suppliedTeamIdentifier string) (valid bool, assignedTeamIdentifier string, err error)
	// GenerateAssetIdentifier generates a new asset identifier and records it as owned by this team
	GenerateAssetIdentifier(ctx context.Context, assetType orgv1.TeamAssetType, assetName string) (string, error)
	// ReuseAssetIdentifier verifies that the supplied identifier is a previously-deleted asset for this
	// team of the relevant type and undeletes it for use with a new asset representing the same logical
	// resource (e.g. a replacement cluster). Returns false if the identifier is not valid, true if it
	// was successfully undeleted.
	ReuseAssetIdentifier(ctx context.Context, assetIdentifier string, assetType orgv1.TeamAssetType, assetName string) (bool, error)
	// MarkAssetDeleted marks a team asset as deleted
	MarkAssetDeleted(ctx context.Context, assetIdentifier string) error
	// ValidateTeamIdentifier checks that the supplied identifier is correct for the team
	ValidateTeamIdentifier(ctx context.Context, teamIdentifier string) (bool, error)
}

type teamAssetsImpl struct {
	team           string
	teamIdentifier string
	teams          Teams
	persist        persistence.Interface
}

func (t *teamAssetsImpl) lookupTeamIdentifier(ctx context.Context, assign bool) (string, error) {
	if t.teamIdentifier != "" {
		return t.teamIdentifier, nil
	}
	team, err := t.persist.Teams().Get(ctx, t.team)
	if err != nil {
		return "", err
	}
	if team.Identifier == "" && assign {
		var err error
		team.Identifier, err = t.teams.GenerateTeamIdentifier(ctx, team.Name)
		if err != nil {
			return "", err
		}

		err = t.persist.Teams().Update(ctx, team)
		if err != nil {
			return "", fmt.Errorf("Failed to update team %s with new identifier: %w", t.team, err)
		}
	}
	t.teamIdentifier = team.Identifier
	return t.teamIdentifier, nil
}

func (t *teamAssetsImpl) EnsureTeamIdentifier(ctx context.Context, suppliedTeamIdentifier string) (valid bool, assignedTeamIdentifier string, err error) {
	valid = false
	// @step: Check if team has an identifier, and assign one if not
	if assignedTeamIdentifier, err = t.lookupTeamIdentifier(ctx, true); err != nil {
		return
	}
	// @step: Validate any supplied identifier matches the assigned one
	if suppliedTeamIdentifier != "" && suppliedTeamIdentifier != assignedTeamIdentifier {
		return
	}
	// @step: Nothing supplied, or supplied matches assigned, either way, it's good.
	valid = true
	return
}

func (t *teamAssetsImpl) GenerateAssetIdentifier(ctx context.Context, assetType orgv1.TeamAssetType, assetName string) (string, error) {
	teamIdentifier, err := t.lookupTeamIdentifier(ctx, true)
	if err != nil {
		return "", err
	}

	assetIdent := utils.GenerateIdentifier()
	err = t.persist.Teams().RecordAsset(ctx, teamIdentifier, assetIdent, model.TeamAssetType(assetType), assetName)
	if err != nil {
		return "", fmt.Errorf("Failed to persist new asset identifier to team %s: %w", t.team, err)
	}
	return assetIdent, nil
}

func (t *teamAssetsImpl) ReuseAssetIdentifier(ctx context.Context, assetIdentifier string, assetType orgv1.TeamAssetType, assetName string) (bool, error) {
	teamIdentifier, err := t.lookupTeamIdentifier(ctx, false)
	if err != nil {
		return false, err
	}

	// Valid for re-use if it exists, is deleted, and previously referred to the same type of asset.
	asset, err := t.persist.Teams().GetAsset(ctx, teamIdentifier, assetIdentifier)
	if err != nil {
		if t.persist.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if asset.DeletedAt == nil || asset.AssetType != model.TeamAssetType(assetType) {
		return false, nil
	}
	if err := t.persist.Teams().MarkAssetUndeleted(ctx, teamIdentifier, assetIdentifier, assetName); err != nil {
		return false, err
	}
	return true, nil
}

func (t *teamAssetsImpl) MarkAssetDeleted(ctx context.Context, assetIdentifier string) error {
	teamIdentifier, err := t.lookupTeamIdentifier(ctx, false)
	if err != nil {
		return err
	}

	return t.persist.Teams().MarkAssetDeleted(ctx, teamIdentifier, assetIdentifier)
}

func (t *teamAssetsImpl) ValidateTeamIdentifier(ctx context.Context, teamIdentifier string) (bool, error) {
	correctTeamIdentifier, err := t.lookupTeamIdentifier(ctx, false)
	if err != nil {
		return false, err
	}
	return teamIdentifier == correctTeamIdentifier, nil
}
