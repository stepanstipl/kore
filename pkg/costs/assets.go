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

package costs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/utils/validation"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Assets provides management over which assets/resources are known to kore, and the known costs
// for those assets provided by a cost provider
type Assets interface {
	ListAssets(ctx context.Context, filters ...persistence.TeamAssetFilterFunc) ([]costsv1.CostAsset, error)
	ListCosts(ctx context.Context, filters ...persistence.TeamAssetFilterFunc) (*costsv1.AssetCostList, error)
	OverallCostsSummary(ctx context.Context, from time.Time, to time.Time, filters ...persistence.TeamAssetFilterFunc) (*costsv1.OverallCostSummary, error)
	TeamCostsSummary(ctx context.Context, teamIdentifier string, from time.Time, to time.Time, filters ...persistence.TeamAssetFilterFunc) (*costsv1.TeamCostSummary, error)
	StoreAssetCosts(ctx context.Context, costs *costsv1.AssetCostList) error
}

// NewAssets creates a new assets implementation
func NewAssets(persistence persistence.TeamAssets, getKoreIdentifier func() string) Assets {
	return &assetsImpl{
		getKoreIdentifier,
		persistence,
	}
}

type assetsImpl struct {
	getKoreIdentifier func() string
	persistence       persistence.TeamAssets
}

func (a *assetsImpl) ListAssets(ctx context.Context, filters ...persistence.TeamAssetFilterFunc) ([]costsv1.CostAsset, error) {
	assets, err := a.persistence.ListAssets(ctx, filters...)
	if err != nil {
		return nil, err
	}
	results := make([]costsv1.CostAsset, len(assets))
	for ind, asset := range assets {
		results[ind] = a.fromTeamAssetModel(asset)
	}
	return results, nil
}

func (a *assetsImpl) ListCosts(ctx context.Context, filters ...persistence.TeamAssetFilterFunc) (*costsv1.AssetCostList, error) {
	costs, err := a.persistence.ListCosts(ctx, filters...)
	if err != nil {
		return nil, err
	}
	results := &costsv1.AssetCostList{
		Items: make([]costsv1.AssetCost, len(costs)),
	}
	for ind, cost := range costs {
		results.Items[ind] = a.fromTeamAssetCostModel(cost)
	}
	return results, nil
}

func (a *assetsImpl) TeamCostsSummary(ctx context.Context, teamIdentifier string, from time.Time, to time.Time, filters ...persistence.TeamAssetFilterFunc) (*costsv1.TeamCostSummary, error) {
	filters = append(filters, persistence.TeamAssetFilters.WithTeam(teamIdentifier))
	filters = append(filters, persistence.TeamAssetFilters.FromTime(from))
	filters = append(filters, persistence.TeamAssetFilters.ToTime(to))
	costs, err := a.persistence.ListCosts(ctx, filters...)
	if err != nil {
		return nil, err
	}

	summary := &costsv1.TeamCostSummary{
		AssetCosts: []*costsv1.AssetCostSummary{},
		CostSummary: costsv1.CostSummary{
			StartTime: metav1.NewTime(from),
			EndTime:   metav1.NewTime(to),
			Cost:      0,
		},
		TeamIdentifier: teamIdentifier,
	}
	assetSummaries := map[string]*costsv1.AssetCostSummary{}

	for _, cost := range costs {
		if summary.TeamName == "" {
			summary.TeamName = cost.Team.TeamName
		}
		if assetSummaries[cost.AssetIdentifier] == nil {
			assetSummaries[cost.AssetIdentifier] = &costsv1.AssetCostSummary{
				AssetIdentifier: cost.AssetIdentifier,
				TeamIdentifier:  cost.TeamIdentifier,
				AssetName:       cost.Asset.AssetName,
				Provider:        cost.Asset.Provider,
				AssetType:       string(cost.Asset.AssetType),
				CostSummary: costsv1.CostSummary{
					StartTime: metav1.NewTime(cost.UsageStartTime),
					EndTime:   metav1.NewTime(cost.UsageEndTime),
					Cost:      0,
				},
				Details: []*costsv1.AssetCost{},
			}
			summary.AssetCosts = append(summary.AssetCosts, assetSummaries[cost.AssetIdentifier])
		}
		asset := assetSummaries[cost.AssetIdentifier]
		asset.Details = a.summariseAssetCosts(asset.Details, cost)
		asset.Cost += cost.Cost
		if asset.StartTime.Time.After(cost.UsageStartTime) {
			asset.StartTime = metav1.NewTime(cost.UsageStartTime)
		}
		if asset.EndTime.Time.Before(cost.UsageEndTime) {
			asset.EndTime = metav1.NewTime(cost.UsageEndTime)
		}
		summary.Cost += cost.Cost
	}
	return summary, nil
}

// summariseAssetCosts is used to group together all usage of the same type for the same resource into a single line item
func (a *assetsImpl) summariseAssetCosts(details []*costsv1.AssetCost, cost *model.TeamAssetCost) []*costsv1.AssetCost {
	// check if usage type already in list
	var detailToUpdate *costsv1.AssetCost = nil
	for _, detail := range details {
		if detail.Description == cost.Description {
			detailToUpdate = detail
			break
		}
	}
	if detailToUpdate == nil {
		c := a.fromTeamAssetCostModel(cost)
		return append(details, &c)
	}
	if detailToUpdate.UsageStartTime.Time.After(cost.UsageStartTime) {
		detailToUpdate.UsageStartTime = metav1.NewTime(cost.UsageStartTime)
	}
	if detailToUpdate.UsageEndTime.Time.Before(cost.UsageEndTime) {
		detailToUpdate.UsageEndTime = metav1.NewTime(cost.UsageEndTime)
	}
	detailToUpdate.Cost += cost.Cost
	amt, err := strconv.ParseFloat(detailToUpdate.UsageAmount, 64)
	if err != nil {
		// Just use the usage amount from this row if we can't parse what's in there, shouldn't happen:
		detailToUpdate.UsageAmount = fmt.Sprintf("%f", cost.UsageAmount)
	} else {
		detailToUpdate.UsageAmount = fmt.Sprintf("%f", amt+cost.UsageAmount)
	}
	return details
}

func (a *assetsImpl) OverallCostsSummary(ctx context.Context, from time.Time, to time.Time, filters ...persistence.TeamAssetFilterFunc) (*costsv1.OverallCostSummary, error) {
	filters = append(filters, persistence.TeamAssetFilters.FromTime(from))
	filters = append(filters, persistence.TeamAssetFilters.ToTime(to))
	costs, err := a.persistence.ListCosts(ctx, filters...)
	if err != nil {
		return nil, err
	}

	summary := &costsv1.OverallCostSummary{
		TeamCosts: []*costsv1.TeamCostSummary{},
		CostSummary: costsv1.CostSummary{
			StartTime: metav1.NewTime(from),
			EndTime:   metav1.NewTime(to),
			Cost:      0,
		},
	}
	teamSummaries := map[string]*costsv1.TeamCostSummary{}

	for _, cost := range costs {
		if teamSummaries[cost.TeamIdentifier] == nil {
			teamSummaries[cost.TeamIdentifier] = &costsv1.TeamCostSummary{
				TeamIdentifier: cost.TeamIdentifier,
				TeamName:       cost.Team.TeamName,
				CostSummary: costsv1.CostSummary{
					StartTime: metav1.NewTime(cost.UsageStartTime),
					EndTime:   metav1.NewTime(cost.UsageEndTime),
					Cost:      0,
				},
			}
			summary.TeamCosts = append(summary.TeamCosts, teamSummaries[cost.TeamIdentifier])
		}
		team := teamSummaries[cost.TeamIdentifier]
		team.Cost += cost.Cost
		if team.StartTime.Time.After(cost.UsageStartTime) {
			team.StartTime = metav1.NewTime(cost.UsageStartTime)
		}
		if team.EndTime.Time.Before(cost.UsageEndTime) {
			team.EndTime = metav1.NewTime(cost.UsageEndTime)
		}
		summary.Cost += cost.Cost
	}
	return summary, nil
}

func (a *assetsImpl) StoreAssetCosts(ctx context.Context, costs *costsv1.AssetCostList) error {
	// parse all first in case there are any errors
	modelCosts := make([]*model.TeamAssetCost, len(costs.Items))
	for ind, cost := range costs.Items {
		var err error
		modelCosts[ind], err = a.toTeamAssetCostModel(cost)
		if err != nil {
			return err
		}
	}
	for ind, modelCost := range modelCosts {
		err := a.persistence.StoreAssetCost(ctx, modelCost)
		if err != nil {
			return fmt.Errorf("error persisting cost index %d: %w", ind, err)
		}
	}
	return nil
}

// fromTeamAssetModel returns a cost asset from the model
func (a *assetsImpl) fromTeamAssetModel(asset *model.TeamAsset) costsv1.CostAsset {
	return costsv1.CostAsset{
		Name:            asset.AssetName,
		AssetIdentifier: asset.AssetIdentifier,
		TeamIdentifier:  asset.TeamIdentifier,
		Provider:        asset.Provider,
		Tags: map[string]string{
			"kore-instance": a.getKoreIdentifier(),
			"kore-team":     asset.TeamIdentifier,
			"kore-cluster":  asset.AssetIdentifier,
		},
	}
}

// fromTeamAssetCostModel returns a asset cost from the model
func (a *assetsImpl) fromTeamAssetCostModel(cost *model.TeamAssetCost) costsv1.AssetCost {
	return costsv1.AssetCost{
		AssetIdentifier: cost.AssetIdentifier,
		TeamIdentifier:  cost.TeamIdentifier,
		Provider:        cost.Provider,
		UsageType:       cost.UsageType,
		Description:     cost.Description,
		UsageStartTime:  metav1.NewTime(cost.UsageStartTime),
		UsageEndTime:    metav1.NewTime(cost.UsageEndTime),
		UsageAmount:     fmt.Sprintf("%f", cost.UsageAmount),
		UsageUnit:       cost.UsageUnit,
		Cost:            cost.Cost,
	}
}

// toTeamAssetCostModel returns a asset cost model from the API type
func (a *assetsImpl) toTeamAssetCostModel(cost costsv1.AssetCost) (*model.TeamAssetCost, error) {
	usageAmount, err := strconv.ParseFloat(cost.UsageAmount, 64)
	if err != nil {
		return nil, validation.NewError("invalid usage amount").WithFieldError("usageAmount", validation.InvalidValue, "cannot parse 'usageAmount' into float")
	}
	return &model.TeamAssetCost{
		AssetIdentifier: cost.AssetIdentifier,
		TeamIdentifier:  cost.TeamIdentifier,
		Provider:        cost.Provider,
		UsageType:       cost.UsageType,
		Description:     cost.Description,
		UsageStartTime:  cost.UsageStartTime.Time,
		UsageEndTime:    cost.UsageEndTime.Time,
		UsageAmount:     usageAmount,
		UsageUnit:       cost.UsageUnit,
		Cost:            cost.Cost,
	}, nil
}
