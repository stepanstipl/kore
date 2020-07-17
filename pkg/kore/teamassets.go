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
	// GenerateAssetIdentifier generates a new asset identifier and records it as owned by this team
	GenerateAssetIdentifier(ctx context.Context, assetType orgv1.TeamAssetType, assetName string) (string, error)
}

type teamAssetsImpl struct {
	team         string
	teams        Teams
	teamsPersist persistence.Teams
}

func (t *teamAssetsImpl) GenerateAssetIdentifier(ctx context.Context, assetType orgv1.TeamAssetType, assetName string) (string, error) {
	team, err := t.teamsPersist.Get(ctx, t.team)
	if err != nil {
		return "", err
	}

	if err := t.ensureTeamHasIdentifier(ctx, team); err != nil {
		return "", err
	}

	assetIdent := utils.GenerateIdentifier()
	err = t.teamsPersist.RecordAsset(ctx, team.Identifier, assetIdent, model.TeamAssetType(assetType), assetName)
	if err != nil {
		return "", fmt.Errorf("Failed to persist new asset identifier to team %s: %w", t.team, err)
	}
	return assetIdent, nil
}

func (t *teamAssetsImpl) ensureTeamHasIdentifier(ctx context.Context, team *model.Team) error {
	// As some teams may pre-date the assignment of team identifiers, ensure the team has
	// an identifier assigned and assign one if not.
	if team.Identifier != "" {
		return nil
	}

	var err error
	team.Identifier, err = t.teams.GenerateTeamIdentifier(ctx, team.Name)
	if err != nil {
		return err
	}

	err = t.teamsPersist.Update(ctx, team)
	if err != nil {
		return fmt.Errorf("Failed to update team %s with new identifier: %w", t.team, err)
	}

	return nil
}
