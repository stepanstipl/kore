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
	"bytes"
	"encoding/json"
	"fmt"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
)

type nodePool struct {
	Name        string
	Size        int64
	MinSize     int64
	MaxSize     int64
	DiskSize    int64
	MachineType string
	Spot        bool
}

type planParamNames struct {
	NodePools   string
	Name        string
	Size        string
	MinSize     string
	MaxSize     string
	DiskSize    string
	MachineType string
	Spot        string
}

func parsePlanConfig(planSpec *configv1.PlanSpec) (map[string]interface{}, error) {
	planConfiguration := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(planSpec.Configuration.Raw)).Decode(&planConfiguration); err != nil {
		return nil, fmt.Errorf("failed to parse plan configuration values: %s", err)
	}
	return planConfiguration, nil
}

func getNodePools(provider string, planConfig map[string]interface{}) ([]nodePool, error) {
	paramNames, err := getPlanParamNames(provider)
	if err != nil {
		return nil, err
	}
	if planConfig[paramNames.NodePools] == nil {
		return nil, nil
	}
	pools := planConfig[paramNames.NodePools].([]interface{})
	nodePools := make([]nodePool, len(pools))
	for i, p := range pools {
		pool := p.(map[string]interface{})
		nodePools[i] = nodePool{
			Name:        pool[paramNames.Name].(string),
			Size:        int64(pool[paramNames.Size].(float64)),
			MinSize:     int64(pool[paramNames.MinSize].(float64)),
			MaxSize:     int64(pool[paramNames.MaxSize].(float64)),
			DiskSize:    int64(pool[paramNames.DiskSize].(float64)),
			MachineType: pool[paramNames.MachineType].(string),
		}
		if paramNames.Spot != "" {
			if pool[paramNames.Spot] != nil {
				nodePools[i].Spot = pool[paramNames.Spot].(bool)
			} else {
				nodePools[i].Spot = false
			}
		}
	}
	return nodePools, nil
}

func getPlanParamNames(provider string) (planParamNames, error) {
	switch provider {
	case providerGCP:
		return planParamNames{
			NodePools:   "nodePools",
			Name:        "name",
			MinSize:     "minSize",
			MaxSize:     "maxSize",
			Size:        "size",
			DiskSize:    "diskSize",
			MachineType: "machineType",
			Spot:        "preemptible",
		}, nil
	case providerAWS:
		return planParamNames{
			NodePools:   "nodeGroups",
			Name:        "name",
			MinSize:     "minSize",
			MaxSize:     "maxSize",
			Size:        "desiredSize",
			DiskSize:    "diskSize",
			MachineType: "instanceType",
			Spot:        "",
		}, nil
	}
	return planParamNames{}, fmt.Errorf("cannot determine plan parameter names for provider %s", provider)
}
