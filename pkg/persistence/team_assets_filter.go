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

import "time"

var (
	// TeamAssetFilters provides a set of useful filters for cost listing
	TeamAssetFilters TeamAssetFilterFuncs
)

// TeamAssetFilterFunc is a function to filter the cost list
type TeamAssetFilterFunc func(*TeamAssetListOptions)

// TeamAssetListOptions is used to specify the restrictions on the returned assets or costs
type TeamAssetListOptions struct {
	From            *time.Time
	To              *time.Time
	BillingMonth    *uint8
	BillingYear     *uint16
	Provider        string
	Account         string
	TeamIdentifier  string
	AssetIdentifier string
	WithDeleted     bool
}

// NewTeamAssetListOptions returns a cost list options
func NewTeamAssetListOptions() *TeamAssetListOptions {
	return &TeamAssetListOptions{
		WithDeleted: false,
	}
}

// ApplyTeamAssetListOptions is responsible for applying the terms
func ApplyTeamAssetListOptions(v ...TeamAssetFilterFunc) *TeamAssetListOptions {
	o := NewTeamAssetListOptions()

	for _, x := range v {
		x(o)
	}

	return o
}

// TeamAssetFilterFuncs provides options for listing assets and asset costs
type TeamAssetFilterFuncs struct{}

// WithTeam limits the returned costs to those where the usage is for assets owned by the specified team identifier
func (q TeamAssetFilterFuncs) WithTeam(teamIdentifier string) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.TeamIdentifier = teamIdentifier
	}
}

// WithAsset limits the returned costs to those where the usage is for the asset specified
func (q TeamAssetFilterFuncs) WithAsset(assetIdentifier string) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.AssetIdentifier = assetIdentifier
	}
}

// FromTime limits the returned costs to those where the usage started after the specified time
func (q TeamAssetFilterFuncs) FromTime(time time.Time) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.From = &time
	}
}

// ToTime limits the returned costs to those where the usage started before the specified time
func (q TeamAssetFilterFuncs) ToTime(time time.Time) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.To = &time
	}
}

// WithBillingPeriod limits the returned costs to those accrued within the specified year/month
// billing period. month = 1 (Jan) to 12 (Dec), year = 4-digit year (e.g. 2020)
func (q TeamAssetFilterFuncs) WithBillingPeriod(year uint16, month uint8) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.BillingMonth = &month
		o.BillingYear = &year
	}
}

// WithProvider limits the returned costs to those accrued with the specified cloud provider (e.g. gcp, aws, azure)
func (q TeamAssetFilterFuncs) WithProvider(provider string) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.Provider = provider
	}
}

// WithAccount limits the returned costs to those accrued within the specified cloud provider account (e.g. AWS account ID, GCP project name)
func (q TeamAssetFilterFuncs) WithAccount(account string) TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.Account = account
	}
}

// WithDeleted includes rows with DeletedAt set in the results
func (q TeamAssetFilterFuncs) WithDeleted() TeamAssetFilterFunc {
	return func(o *TeamAssetListOptions) {
		o.WithDeleted = true
	}
}
