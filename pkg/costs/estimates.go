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
	"errors"
	"fmt"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/utils/validation"
)

// Estimates allows requesting of estimation of the running cost of a plan
type Estimates interface {
	// GetClusterEstimate gives an estimate of the running costs for a cluster based on the provided plan
	GetClusterEstimate(planSpec *configv1.PlanSpec) (*costsv1.CostEstimate, error)
	// GetServiceEstimate gives an estimate of the running costs for a service based on the provided plan
	GetServiceEstimate(planSpec *servicesv1.ServicePlanSpec) (*costsv1.CostEstimate, error)
}

// NewEstimates creates a new instance of the estimates API
func NewEstimates(metadata Metadata) Estimates {
	return &estimatesImpl{
		metadata,
	}
}

type estimatesImpl struct {
	metadata Metadata
}

const (
	cloudGCP      = "gcp"
	cloudAWS      = "aws"
	cloudAzure    = "azure"
	providerGCP   = "GKE"
	providerAWS   = "EKS"
	providerAzure = "AKS"
)

func (e *estimatesImpl) GetClusterEstimate(planSpec *configv1.PlanSpec) (*costsv1.CostEstimate, error) {
	// Load costing info for the provider in question

	cloud := getCloudForClusterProvider(planSpec.Kind)
	if cloud == "" {
		// Cost estimation not supported as it's not provided by a cloud
		return nil, validation.NewError("plan not valid").WithFieldErrorf("kind", validation.InvalidValue, "cannot determine cloud provider for cluster provider %s", planSpec.Kind)
	}

	// Determine region
	planConfiguration, err := parsePlanConfig(planSpec)
	if err != nil {
		return nil, err
	}
	region, ok := planConfiguration["region"].(string)
	if !ok || region == "" {
		// Cost estimation not supported until a region is selected
		return nil, validation.NewError("plan not valid").WithFieldError("region", validation.Required, "region required to produce estimate")
	}

	estimate := &costsv1.CostEstimate{}

	// Start with control plane cost
	controlPlaneCost, err := e.metadata.KubernetesControlPlaneCost(cloud, region)
	if err != nil {
		return nil, err
	}
	estimate.CostElements = append(estimate.CostElements, costsv1.CostEstimateElement{
		Name:        "Control Plane",
		MinCost:     controlPlaneCost,
		MaxCost:     controlPlaneCost,
		TypicalCost: controlPlaneCost,
	})

	// Add cost for the load balancer we always deploy for the auth proxy
	exposedServiceCost, err := e.metadata.KubernetesExposedServiceCost(cloud, region)
	if err != nil {
		return nil, err
	}
	estimate.CostElements = append(estimate.CostElements, costsv1.CostEstimateElement{
		Name:        "Kore Authentication Load Balancer",
		MinCost:     exposedServiceCost,
		MaxCost:     exposedServiceCost,
		TypicalCost: exposedServiceCost,
	})

	// Add node pool costs
	nodePools, err := getNodePools(planSpec.Kind, planConfiguration)
	if err != nil {
		return nil, err
	}

	zoneMultiplier := int64(1)
	if planSpec.Kind == providerGCP {
		// GKE is hard-coded currently to deploy for all zones in a region, so get the zones for the region
		// and multiply the node pool size by that to get the estimate.
		regionAZs, err := e.metadata.RegionZones(cloud, region)
		if err != nil {
			return nil, err
		}
		zoneMultiplier = int64(len(regionAZs))
	}

	for _, nodePool := range nodePools {
		npEstimate, err := e.GetNodePoolEstimate(cloud, region, nodePool, zoneMultiplier)
		if err != nil {
			return nil, err
		}
		estimate.CostElements = append(estimate.CostElements, *npEstimate)
	}

	// @TODO: HTTP load balancers, other chargeable cluster features like Anthos?
	// @TODO: Non-EKS VPC resources in AWS we create as part of a cluster

	// Create overall summary
	for _, e := range estimate.CostElements {
		estimate.MinCost += e.MinCost
		estimate.MaxCost += e.MaxCost
		estimate.TypicalCost += e.TypicalCost
	}

	return estimate, nil
}

func (e *estimatesImpl) GetNodePoolEstimate(cloud string, region string, nodePool nodePool, zoneMultiplier int64) (*costsv1.CostEstimateElement, error) {
	instType, err := e.metadata.InstanceType(cloud, region, nodePool.MachineType)
	if err != nil {
		return nil, err
	}
	if instType == nil {
		return nil, validation.NewError("plan not valid").
			WithFieldErrorf(fmt.Sprintf("nodePool.%s.instanceType", nodePool.Name), validation.InvalidValue, "no price available for %s in %s - the instance type may not be available in this region or the region may not exist", nodePool.MachineType, region)
	}

	priceType := costsv1.PriceTypeOnDemand
	if nodePool.Spot {
		if cloud == cloudGCP {
			priceType = costsv1.PriceTypePreEmptible
		} else {
			priceType = costsv1.PriceTypeSpot
		}
	}
	nodePrice := instType.Prices[priceType]
	if nodePrice == 0 {
		return nil, validation.NewError("plan not valid").
			WithFieldErrorf(fmt.Sprintf("nodePool.%s.priceType", nodePool.Name), validation.InvalidValue, "no price of type %s available for %s in %s - this price type may not be available in this region", priceType, nodePool.MachineType, region)
	}
	estimate := &costsv1.CostEstimateElement{
		Name:        fmt.Sprintf("Node Pool %s", nodePool.Name),
		TypicalCost: nodePool.Size * nodePrice * zoneMultiplier,
	}
	if nodePool.AutoScale {
		estimate.MinCost = nodePool.MinSize * nodePrice * zoneMultiplier
		estimate.MaxCost = nodePool.MaxSize * nodePrice * zoneMultiplier
	} else {
		estimate.MinCost = estimate.TypicalCost
		estimate.MaxCost = estimate.TypicalCost
	}

	return estimate, nil
}

func (*estimatesImpl) GetServiceEstimate(planSpec *servicesv1.ServicePlanSpec) (*costsv1.CostEstimate, error) {
	return nil, errors.New("Not implemented")
}
