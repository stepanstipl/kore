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

package aks

import (
	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
)

func ConvertAuthorizedMasterNetworks(v []AuthorizedMasterNetwork) []string {
	res := make([]string, 0, len(v))
	for _, e := range v {
		res = append(res, e.Cidr)
	}
	return res
}

func ConvertSSHPublicKeys(v []SSHPublicKey) []string {
	res := make([]string, 0, len(v))
	for _, e := range v {
		res = append(res, string(e))
	}
	return res
}

func ConvertLinuxProfile(v *LinuxProfile) *aksv1alpha1.LinuxProfile {
	if v == nil {
		return nil
	}

	return &aksv1alpha1.LinuxProfile{
		AdminUsername: v.AdminUsername,
		SSHPublicKeys: ConvertSSHPublicKeys(v.SSHPublicKeys),
	}
}

func ConvertWindowsProfile(v *WindowsProfile) *aksv1alpha1.WindowsProfile {
	if v == nil {
		return nil
	}

	return &aksv1alpha1.WindowsProfile{
		AdminUsername: v.AdminUsername,
		AdminPassword: v.AdminPassword,
	}
}

func ConvertNodePools(v []NodePool) []aksv1alpha1.AgentPoolProfile {
	res := make([]aksv1alpha1.AgentPoolProfile, 0, len(v))
	for _, e := range v {
		res = append(res, aksv1alpha1.AgentPoolProfile{
			Name:              e.Name,
			Mode:              e.Mode,
			EnableAutoScaling: e.EnableAutoscaler,
			NodeImageVersion:  e.Version,
			Count:             e.Size,
			MinCount:          e.MinSize,
			MaxCount:          e.MaxSize,
			MaxPods:           e.MaxPodsPerNode,
			VMSize:            e.MachineType,
			OsType:            e.ImageType,
			OsDiskSizeGB:      e.DiskSize,
			NodeLabels:        ConvertLabels(e.Labels),
			NodeTaints:        ConvertTaints(e.Taints),
		})
	}
	return res
}

func ConvertLabels(v map[string]Label) map[string]string {
	res := make(map[string]string, len(v))
	for k, e := range v {
		res[k] = string(e)
	}
	return res
}

func ConvertTaints(v []Taint) []aksv1alpha1.NodeTaint {
	res := make([]aksv1alpha1.NodeTaint, 0, len(v))
	for _, e := range v {
		res = append(res, aksv1alpha1.NodeTaint{
			Key:    e.Key,
			Value:  e.Value,
			Effect: e.Effect,
		})
	}
	return res
}
