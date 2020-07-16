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
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var plans = []configv1.Plan{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aks-development",
			Annotations: map[string]string{
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Provides a development AKS cluster",
			Description: "AKS Development Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "development",
			},
			Configuration: apiextv1.JSON{
				Raw: []byte(`{
					"authorizedMasterNetworks": [
						{
							"name": "default",
							"cidr": "0.0.0.0/0"
						}
					],
					"authProxyAllowedIPs":       ["0.0.0.0/0"],
					"defaultTeamRole":           "view",
					"description":               "AKS Development Cluster",
					"dnsPrefix":                 "kore",
					"domain":                    "default",
					"enableDefaultTrafficBlock": false,
					"enablePodSecurityPolicy":   false,
					"inheritTeamMembers":        true,
					"networkPlugin":             "kubenet",
					"privateClusterEnabled":     false,
					"region":                    "uksouth",
					"nodePools": [
						{
							"name":             "system",
							"mode":             "System",
							"size":             1,
							"minSize":          1,
							"maxSize":          10,
							"machineType":      "Standard_D2_v2",
							"imageType":        "Linux",
							"diskSize":         30,
							"enableAutoscaler": true
						}
					],
					"version":                  "1.16.10"
				}`),
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aks-production",
			Annotations: map[string]string{
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Provides a production AKS cluster",
			Description: "AKS Production Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "production",
			},
			Configuration: apiextv1.JSON{
				Raw: []byte(`{
					"authorizedMasterNetworks": [
						{
							"name": "default",
							"cidr": "0.0.0.0/0"
						}
					],
					"authProxyAllowedIPs":       ["0.0.0.0/0"],
					"defaultTeamRole":           "view",
					"description":               "AKS Production Cluster",
					"dnsPrefix":                 "kore",
					"domain":                    "default",
					"enableDefaultTrafficBlock": false,
					"enablePodSecurityPolicy":   false,
					"inheritTeamMembers":        true,
					"networkPlugin":             "kubenet",
					"privateClusterEnabled":     false,
					"region":                    "uksouth",
					"nodePools": [
						{
							"name":             "system",
							"mode":             "System",
							"size":             3,
							"minSize":          3,
							"maxSize":          10,
							"machineType":      "Standard_D2_v2",
							"imageType":        "Linux",
							"diskSize":         30,
							"enableAutoscaler": true
						}
					],
					"version":                  "1.16.10"
				}`),
			},
		},
	},
}
