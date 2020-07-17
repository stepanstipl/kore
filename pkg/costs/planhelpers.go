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
	"github.com/appvia/kore/pkg/utils/validation"
)

type nodePool struct {
	Name        string
	Size        int64
	MinSize     int64
	MaxSize     int64
	DiskSize    int64
	MachineType string
	Spot        bool
	AutoScale   bool
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
	AutoScale   string
}

func parsePlanConfig(planSpec *configv1.PlanSpec) (map[string]interface{}, error) {
	planConfiguration := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(planSpec.Configuration.Raw)).Decode(&planConfiguration); err != nil {
		return nil, fmt.Errorf("failed to parse plan configuration values: %s", err)
	}
	return planConfiguration, nil
}

func nodePoolValError(npInd int, field string, msg string) error {
	return validation.NewError("plan not valid").WithFieldError(
		fmt.Sprintf("nodePool[%v].%s", npInd, field),
		validation.InvalidValue,
		msg)
}

func parseNodePoolIntProperty(pool map[string]interface{}, paramName string, npInd int, msg string) (int64, error) {
	val, ok := pool[paramName].(float64)
	if !ok {
		return 0, nodePoolValError(npInd, paramName, msg)
	}
	return int64(val), nil
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
		nodePools[i] = nodePool{}
		ok := false
		var err error
		if nodePools[i].Name, ok = pool[paramNames.Name].(string); !ok {
			return nil, nodePoolValError(i, paramNames.Name, "all node pools must have a name to produce estimate")
		}
		if nodePools[i].Size, err = parseNodePoolIntProperty(pool, paramNames.Size, i, "all node pools must have a size to produce estimate"); err != nil {
			return nil, err
		}
		if nodePools[i].DiskSize, err = parseNodePoolIntProperty(pool, paramNames.DiskSize, i, "invalid disk size"); err != nil {
			return nil, err
		}
		if nodePools[i].AutoScale, ok = pool[paramNames.AutoScale].(bool); !ok {
			return nil, nodePoolValError(i, paramNames.AutoScale, "invalid auto-scale")
		}
		if nodePools[i].MachineType, ok = pool[paramNames.MachineType].(string); !ok {
			return nil, nodePoolValError(i, paramNames.MachineType, "invalid machine type")
		}
		if nodePools[i].AutoScale {
			if nodePools[i].MinSize, err = parseNodePoolIntProperty(pool, paramNames.MinSize, i, "node pools must have a min size to produce estimate with auto-scale enabled"); err != nil {
				return nil, err
			}
			if nodePools[i].MaxSize, err = parseNodePoolIntProperty(pool, paramNames.MaxSize, i, "node pools must have a max size to produce estimate with auto-scale enabled"); err != nil {
				return nil, err
			}
		}
		if paramNames.Spot != "" {
			if pool[paramNames.Spot] != nil {
				if nodePools[i].Spot, ok = pool[paramNames.Spot].(bool); !ok {
					return nil, nodePoolValError(i, paramNames.Spot, "invalid spot")
				}
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
			AutoScale:   "enableAutoscaler",
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
			AutoScale:   "enableAutoscaler",
		}, nil
	case providerAzure:
		return planParamNames{
			NodePools:   "nodePools",
			Name:        "name",
			MinSize:     "minSize",
			MaxSize:     "maxSize",
			Size:        "size",
			DiskSize:    "diskSize",
			MachineType: "machineType",
			Spot:        "",
			AutoScale:   "enableAutoscaler",
		}, nil
	}
	return planParamNames{}, fmt.Errorf("cannot determine plan parameter names for provider %s", provider)
}
