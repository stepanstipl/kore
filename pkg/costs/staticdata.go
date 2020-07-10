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
	"encoding/json"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/costs/staticdata"
)

// staticData provides an implementation of the Cloudinfo interface that uses static
// data instead of live data, without any pricing information
type staticData struct {
}

var _ Cloudinfo = &staticData{}

func (s *staticData) Ready() bool {
	return true
}

func (s *staticData) KubernetesRegions(cloud string) ([]costsv1.Continent, error) {
	j := staticdata.Continents(cloud)
	if j == "" {
		return nil, nil
	}
	var continents []costsv1.Continent
	err := json.Unmarshal([]byte(j), &continents)
	if err != nil {
		return nil, err
	}
	return continents, nil
}

func (s *staticData) KubernetesRegionAZs(cloud string, region string) ([]string, error) {
	return nil, nil
}

func (s *staticData) KubernetesInstanceTypes(cloud string, region string) ([]costsv1.InstanceType, error) {
	j := staticdata.NodeTypes(cloud)
	if j == "" {
		return nil, nil
	}
	var res map[string]interface{}
	err := json.Unmarshal([]byte(j), &res)
	if err != nil {
		return nil, err
	}

	result := []costsv1.InstanceType{}
	for _, product := range res["products"].([]interface{}) {
		pm := product.(map[string]interface{})
		if pm["category"].(string) == "" {
			continue
		}
		info := costsv1.InstanceType{
			Name:     pm["type"].(string),
			Category: pm["category"].(string),
			MCpus:    int64(pm["cpusPerVm"].(float64) * 1000),
			Mem:      int64(pm["memPerVm"].(float64) * 1000),
		}
		result = append(result, info)
	}
	return result, nil
}

func (s *staticData) KubernetesInstanceType(cloud string, region string, instanceType string) (*costsv1.InstanceType, error) {
	types, err := s.KubernetesInstanceTypes(cloud, region)
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

func (s *staticData) KubernetesVersions(cloud string, region string) ([]string, error) {
	return nil, nil
}
