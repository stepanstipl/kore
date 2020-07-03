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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CostEstimate defines the result of the cost estimation
// +k8s:openapi-gen=false
type CostEstimate struct {
	// MinCost is the minimum hourly cost estimate in microdollars
	MinCost int64 `json:"minCost,omitempty"`
	// TypicalCost is the expected / likely hourly cost estimate in microdollars
	TypicalCost int64 `json:"typicalCost,omitempty"`
	// MaxCost is the estimated upper limit of the hourly cost in microdollars
	MaxCost int64 `json:"maxCost,omitempty"`
	// CostElements provides details of the different components which make up this cost estimate
	CostElements []CostEstimateElement `json:"costElements,omitempty"`
	// PreparedAt indicates the time this estimate was prepared
	PreparedAt metav1.Time `json:"preparedAt,omitempty"`
}

// CostEstimateElement represents a logical component which has an associated cost
// +k8s:openapi-gen=false
type CostEstimateElement struct {
	// Name is the name of this component
	Name string `json:"name,omitempty"`
	// MinCost is the minimum hourly cost estimate of this component in microdollars
	MinCost int64 `json:"minCost,omitempty"`
	// TypicalCost is the expected / likely hourly cost estimate of this component in microdollars
	TypicalCost int64 `json:"typicalCost,omitempty"`
	// MaxCost is the estimated upper limit of the hourly cost of this component in microdollars
	MaxCost int64 `json:"maxCost,omitempty"`
}

// ContinentList provides the list of continents and regions available for a cloud
// provider, which can further be used to request price details for that provider
// in a specific region.
// +k8s:openapi-gen=false
type ContinentList struct {
	Items []Continent `json:"items"`
}

// Continent is a geographical grouping of regions
// +k8s:openapi-gen=false
type Continent struct {
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}

// Region is a specific cloud provider region
// +k8s:openapi-gen=false
type Region struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// PriceType is the possible types of prices for cloud infrastructure
// +k8s:openapi-gen=false
type PriceType string

// PriceTypeOnDemand is the normal 'rack' price for a piece of infrastructure
const PriceTypeOnDemand PriceType = "OnDemand"

// PriceTypeSpot is the variable price which you may be able to use a piece of infrastructure for
const PriceTypeSpot PriceType = "Spot"

// PriceTypePreEmptible is the fixed discounted price which you can use a piece of infrastructure for subject to availability and early termination
const PriceTypePreEmptible PriceType = "PreEmptible"

// InstanceType is an available compute type from a cloud provider
// +k8s:openapi-gen=false
type InstanceType struct {
	// Category is the classification of this instance type
	Category string `json:"category"`
	// Name is the unique identifier of this instance type
	Name string `json:"name"`
	// Prices gives the price of this instance type in microdollars per hour for the given price type
	Prices map[PriceType]int64 `json:"prices"`
	// MCpus is the number of milliCPUs assigned to this instance type
	MCpus int64 `json:"mCpus"`
	// Mem is the amount of memory, expressed in milli-GiBs, assigned to this instance type
	Mem int64 `json:"mem"`
}

// InstanceTypeList is a group of instance types available for a cloud provider
// +k8s:openapi-gen=false
type InstanceTypeList struct {
	Items []InstanceType `json:"items"`
}
