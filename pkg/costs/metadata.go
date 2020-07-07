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
	"fmt"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
)

// Metadata allows access to cloud service metadats such as instance types and prices
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Metadata
type Metadata interface {
	// Clouds retrieves the list of supported clouds
	Clouds() ([]string, error)
	// MapProviderToCloud maps a provider to a cloud
	MapProviderToCloud(provider string) (string, error)
	// Regions retrieves the list of available regions, organised by continent, for the specified cloud
	Regions(cloud string) (*costsv1.ContinentList, error)
	// RegionZones retrieves the list of available AZs in the given region, for the specified cloud
	RegionZones(cloud string, region string) ([]string, error)
	// InstanceTypes retrieves the list of available instance types for the specified cloud and region
	InstanceTypes(cloud string, region string) (*costsv1.InstanceTypeList, error)
	// InstanceType gets the metadata for a specific selected instance type for the specified cloud and region
	InstanceType(cloud string, region string, instanceType string) (*costsv1.InstanceType, error)
	// KubernetesVersions retrieves the list of supported kubernetes versions for the specified cloud and region
	KubernetesVersions(cloud string, region string) ([]string, error)
	// KubernetesControlPlanCost retrieves the price in microdollars per hour of a Kubernetes
	// control plane in the specific cloud and region
	KubernetesControlPlaneCost(cloud string, region string) (int64, error)
	// KubernetesExposedServiceCost retrieves the price in microdollars per hour of an exposed service (i.e.
	// HTTP load balancer)
	KubernetesExposedServiceCost(cloud string, region string) (int64, error)
}

// NewMetadata creates a new instance of the metadata API
func NewMetadata(cloudinfo Cloudinfo) Metadata {
	return &metadataImpl{
		cloudinfo,
	}
}

type metadataImpl struct {
	cloudinfo Cloudinfo
}

func (m *metadataImpl) Clouds() ([]string, error) {
	return []string{cloudGCP, cloudAWS, cloudAzure}, nil
}

func (m *metadataImpl) MapProviderToCloud(provider string) (string, error) {
	cloud := getCloudForClusterProvider(provider)
	if cloud == "" {
		return "", fmt.Errorf("unknown Kubernetes provider %s, cannot determine cloud provider", provider)
	}
	return cloud, nil
}

func (m *metadataImpl) Regions(cloud string) (*costsv1.ContinentList, error) {
	continents, err := m.cloudinfo.KubernetesRegions(cloud)
	if err != nil {
		return nil, err
	}
	if continents == nil {
		return nil, nil
	}
	result := &costsv1.ContinentList{}
	result.Items = append(result.Items, continents...)
	return result, nil
}

func (m *metadataImpl) RegionZones(cloud string, region string) ([]string, error) {
	return m.cloudinfo.KubernetesRegionAZs(cloud, region)
}

func (m *metadataImpl) InstanceTypes(cloud string, region string) (*costsv1.InstanceTypeList, error) {
	instanceTypes, err := m.cloudinfo.KubernetesInstanceTypes(cloud, region)
	if err != nil {
		return nil, err
	}
	if instanceTypes == nil {
		return nil, nil
	}
	result := &costsv1.InstanceTypeList{}
	result.Items = append(result.Items, instanceTypes...)
	return result, nil
}

func (m *metadataImpl) InstanceType(cloud string, region string, instanceType string) (*costsv1.InstanceType, error) {
	return m.cloudinfo.KubernetesInstanceType(cloud, region, instanceType)
}

func (m *metadataImpl) KubernetesVersions(cloud string, region string) ([]string, error) {
	return m.cloudinfo.KubernetesVersions(cloud, region)
}

func (m *metadataImpl) KubernetesControlPlaneCost(cloud string, region string) (int64, error) {
	// @TODO: Determine this from the providers. For now, AWS and GCP both charge 10c/hr worldwide.
	switch cloud {
	case cloudGCP:
		return 100000, nil
	case cloudAWS:
		return 100000, nil
	}
	return 0, nil
}

func (m *metadataImpl) KubernetesExposedServiceCost(cloud string, region string) (int64, error) {
	// @TODO: Determine this from the providers. For now, just return a typicalish load
	// balancer cost of ~2.5c/hr
	return 25000, nil
}
