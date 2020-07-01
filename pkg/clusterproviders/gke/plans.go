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

package gke

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var plans = []configv1.Plan{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gke-development",
			Annotations: map[string]string{
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Provides a development cluster within GKE",
			Description: "GKE Development Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "dev",
				kore.Label("kind"):        Kind,
				kore.Label("plural"):      "gkes",
			},
			Configuration: apiextv1.JSON{
				Raw: []byte(`{
					"authorizedMasterNetworks": [
						{
							"name": "default",
							"cidr": "0.0.0.0/0"
						}
					],
					"authProxyAllowedIPs":           ["0.0.0.0/0"],
					"defaultTeamRole":               "view",
					"description":                   "gke-development cluster",
					"domain":                        "default",
					"enableDefaultTrafficBlock":     false,
					"enableHTTPLoadBalancer":        true,
					"enableHorizontalPodAutoscaler": true,
					"enableIstio":                   false,
					"enablePrivateEndpoint":         false,
					"enablePrivateNetwork":          false,
					"enableShieldedNodes":           true,
					"enableStackDriverLogging":      true,
					"enableStackDriverMetrics":      true,
					"inheritTeamMembers":            true,
					"maintenanceWindow":             "03:00",
					"region":                        "europe-west2",
					"version":                       "",
					"releaseChannel":                "REGULAR",
					"nodePools": [
						{
							"name":                          "compute",
							"enableAutoupgrade":             true,
							"version":                       "",
							"enableAutoscaler":              true,
							"enableAutorepair":              true,
							"minSize":                       1,
							"maxSize":                       10,
							"size":                          1,
							"maxPodsPerNode":                110,
							"machineType":                   "n1-standard-2",
							"imageType":                     "COS",
							"diskSize":                      100,
							"preemptible":                   false
						}
					]
				}`),
			},
		},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gke-production",
			Annotations: map[string]string{
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Provides a production cluster within GKE",
			Description: "GKE Production Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "production",
				kore.Label("kind"):        Kind,
				kore.Label("plural"):      "gkes",
			},
			Configuration: apiextv1.JSON{
				Raw: []byte(`{
					"authorizedMasterNetworks": [
						{
							"name": "default",
							"cidr": "0.0.0.0/0"
						}
					],
					"authProxyAllowedIPs":           ["0.0.0.0/0"],
					"defaultTeamRole":               "view",
					"description":                   "gke-production cluster",
					"domain":                        "default",
					"enableDefaultTrafficBlock":     false,
					"enableHTTPLoadBalancer":        true,
					"enableHorizontalPodAutoscaler": true,
					"enableIstio":                   false,
					"enablePrivateEndpoint":         false,
					"enablePrivateNetwork":          false,
					"enableShieldedNodes":           true,
					"enableStackDriverLogging":      true,
					"enableStackDriverMetrics":      true,
					"inheritTeamMembers":            true,
					"maintenanceWindow":             "03:00",
					"region":                        "europe-west2",
					"version":                       "",
					"releaseChannel":                "REGULAR",
					"nodePools": [
						{
							"name":                          "compute",
							"enableAutoupgrade":             true,
							"version":                       "",
							"enableAutoscaler":              true,
							"enableAutorepair":              true,
							"minSize":                       1,
							"maxSize":                       10,
							"size":                          2,
							"maxPodsPerNode":                110,
							"machineType":                   "n1-standard-2",
							"imageType":                     "COS",
							"diskSize":                      100,
							"preemptible":                   false
						}
					]
				}`),
			},
		},
	},
}
