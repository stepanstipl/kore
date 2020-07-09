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
	"strings"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"

	"github.com/go-resty/resty/v2"
)

// Cloudinfo provides an abstracted interface to the cloudinfo service
type Cloudinfo interface {
	KubernetesRegions(cloud string) ([]costsv1.Continent, error)
	KubernetesRegionAZs(cloud string, region string) ([]string, error)
	KubernetesInstanceTypes(cloud string, region string) ([]costsv1.InstanceType, error)
	KubernetesInstanceType(cloud string, region string, instanceType string) (*costsv1.InstanceType, error)
	KubernetesVersions(cloud string, region string) ([]string, error)
}

var _ Cloudinfo = &cloudinfoImpl{}

type cloudinfoImpl struct {
	url string
}

// NewCloudInfo returns a new cloudinfo interface
func NewCloudInfo(url string) Cloudinfo {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return &cloudinfoImpl{
		url: url,
	}
}

func (c *cloudinfoImpl) providerServiceURL(cloud string) (string, error) {
	cloudinfCloud, service, err := mapCloud(cloud)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%sapi/v1/providers/%s/services/%s", c.url, cloudinfCloud, service), nil
}

func (c *cloudinfoImpl) providerServiceRegionURL(cloud string, region string) (string, error) {
	cloudinfCloud, service, err := mapCloud(cloud)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%sapi/v1/providers/%s/services/%s/regions/%s", c.url, cloudinfCloud, service, region), nil
}

func (c *cloudinfoImpl) KubernetesRegions(cloud string) ([]costsv1.Continent, error) {
	client := resty.New()
	url, err := c.providerServiceURL(cloud)
	if err != nil {
		return nil, err
	}

	resp, err := client.R().
		SetResult([]costsv1.Continent{}).
		Get(url + "/continents")
	if err != nil {
		return nil, err
	}

	continents := resp.Result().(*[]costsv1.Continent)
	return *continents, nil
}

func (c *cloudinfoImpl) KubernetesRegionAZs(cloud string, region string) ([]string, error) {
	client := resty.New()
	url, err := c.providerServiceRegionURL(cloud, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.R().
		SetResult(map[string]interface{}{}).
		Get(url)
	if err != nil {
		return nil, err
	}

	rawResult := resp.Result()
	if rawResult == nil {
		return nil, nil
	}
	zones := (*rawResult.(*map[string]interface{}))["zones"]
	if zones == nil {
		// TODO: This suggests an error if the response itself was non-nil, this SHOULD
		// have something...
		return nil, nil
	}
	zoneList := zones.([]interface{})
	ret := make([]string, len(zoneList))
	for i, z := range zoneList {
		ret[i] = z.(string)
	}
	return ret, nil
}

func (c *cloudinfoImpl) KubernetesInstanceTypes(cloud string, region string) ([]costsv1.InstanceType, error) {
	client := resty.New()
	url, err := c.providerServiceRegionURL(cloud, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.R().
		SetResult(map[string]interface{}{}).
		Get(url + "/products")
	if err != nil {
		return nil, err
	}

	rawResult := resp.Result()
	if rawResult == nil {
		return nil, nil
	}
	products := (*rawResult.(*map[string]interface{}))["products"]
	if products == nil {
		// TODO: This suggests an error if the response itself was non-nil, this SHOULD
		// have something...
		return nil, nil
	}

	result := []costsv1.InstanceType{}
	for _, product := range products.([]interface{}) {
		pm := product.(map[string]interface{})
		// Cloudinfo returns a bunch of 'eks control plan' compute products for EKS, all of which are identical
		// and/or meaningless - filter them out by now (they have an empty category)
		if pm["category"].(string) == "" {
			continue
		}
		info := costsv1.InstanceType{
			Name:     pm["type"].(string),
			Category: pm["category"].(string),
			MCpus:    int64(pm["cpusPerVm"].(float64) * 1000),
			Mem:      int64(pm["memPerVm"].(float64) * 1000),
			Prices: map[costsv1.PriceType]int64{
				costsv1.PriceTypeOnDemand: parsePriceToMicrodollars(pm["onDemandPrice"]),
			},
		}
		if pm["spotPrice"] != nil && pm["spotPrice"].([]interface{}) != nil {
			// for now, simply take the /first/ spot price. if GCP, report as pre-emptible price,
			// otherwise, report as spot price.
			spotPriceType := costsv1.PriceTypeSpot
			if cloud == "gcp" {
				spotPriceType = costsv1.PriceTypePreEmptible
			}
			info.Prices[spotPriceType] = parsePriceToMicrodollars(pm["spotPrice"].([]interface{})[0].(map[string]interface{})["price"])
		}
		result = append(result, info)
	}
	return result, nil
}

func (c *cloudinfoImpl) KubernetesInstanceType(cloud string, region string, instanceType string) (*costsv1.InstanceType, error) {
	// For now, just call the list and find the instance type. Pretty inefficient but probably
	// good enough.
	types, err := c.KubernetesInstanceTypes(cloud, region)
	if err != nil {
		return nil, err
	}
	for _, t := range types {
		if t.Name == instanceType {
			return &t, nil
		}
	}
	return nil, nil
}

func (c *cloudinfoImpl) KubernetesVersions(cloud string, region string) ([]string, error) {
	client := resty.New()
	url, err := c.providerServiceRegionURL(cloud, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.R().
		SetResult([]interface{}{}).
		Get(url + "/versions")
	if err != nil {
		return nil, err
	}

	versions := resp.Result().(*[]interface{})
	if versions == nil || len(*versions) == 0 {
		return nil, nil
	}

	// For now, just return the versions for the first reported AZ
	versionList := (*versions)[0].(map[string]interface{})["versions"].([]interface{})
	ret := make([]string, len(versionList))
	for i, v := range versionList {
		ret[i] = v.(string)
	}
	return ret, nil
}

// mapCloud translates our cloud names into cloudinfo cloud names and kubernetes service names
func mapCloud(cloud string) (string, string, error) {
	switch cloud {
	case cloudGCP:
		return "google", "gke", nil
	case cloudAWS:
		return "amazon", "eks", nil
	case cloudAzure:
		return "azure", "aks", nil
	}
	return "", "", fmt.Errorf("Can't map provided cloud %s to cloudinfo cloud name", cloud)
}

func parsePriceToMicrodollars(val interface{}) int64 {
	dollars := val.(float64)
	return int64(dollars * 1000000)
}
