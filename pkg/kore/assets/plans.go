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

package assets

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Label returns a kore label
func Label(name string) string {
	return "kore.appvia.io/" + name
}

// GetDefaultPlans returns a collection of plans for the resources
func GetDefaultPlans() []*configv1.Plan {
	return []*configv1.Plan{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kore",
				Annotations: map[string]string{
					"kore.appvia.io/system":   "true",
					"kore.appvia.io/readonly": "true",
				},
			},
			Spec: configv1.PlanSpec{
				Kind:        "Kore",
				Summary:     "Default cluster plan for Kore",
				Description: "Default cluster plan for Kore",
				Configuration: apiextv1.JSON{
					Raw: []byte(`{}`),
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gke-development",
				Annotations: map[string]string{
					"kore.appvia.io/readonly": "true",
				},
			},
			Spec: configv1.PlanSpec{
				Kind:        "GKE",
				Summary:     "Provides a development cluster within GKE",
				Description: "GKE Development Cluster",
				Labels: map[string]string{
					Label("environment"): "dev",
					Label("kind"):        "GKE",
					Label("plural"):      "gkes",
				},
				Configuration: apiextv1.JSON{
					Raw: []byte(`
						{
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
						}`,
					),
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
				Kind:        "GKE",
				Summary:     "Provides a production cluster within GKE",
				Description: "GKE Production Cluster",
				Labels: map[string]string{
					Label("environment"): "production",
					Label("kind"):        "GKE",
					Label("plural"):      "gkes",
				},
				Configuration: apiextv1.JSON{
					Raw: []byte(`
						{
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
						}`,
					),
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "eks-development",
				Annotations: map[string]string{
					"kore.appvia.io/readonly": "true",
				},
			},
			Spec: configv1.PlanSpec{
				Kind:        "EKS",
				Summary:     "Provides a development cluster within EKS",
				Description: "EKS Development Cluster",
				Labels: map[string]string{
					Label("environment"): "development",
					Label("kind"):        "EKS",
					Label("plural"):      "ekss",
				},
				Configuration: apiextv1.JSON{
					Raw: []byte(`{
						"authProxyAllowedIPs": [
							"0.0.0.0/0"
						],
						"defaultTeamRole": "view",
						"description": "eks-development cluster",
						"domain": "default",
						"enableDefaultTrafficBlock": false,
						"inheritTeamMembers": true,
						"privateIPV4Cidr": "10.0.0.0/16",
						"region": "eu-west-2",
						"version": "1.15",
						"nodeGroups": [
							{
								"name": "default",
								"instanceType": "t3.medium",
								"diskSize": 10,
								"enableAutoscaler": true,
								"minSize": 1,
								"desiredSize": 1,
								"maxSize": 10
							}
						]
					}`),
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "eks-production",
				Annotations: map[string]string{
					"kore.appvia.io/readonly": "true",
				},
			},
			Spec: configv1.PlanSpec{
				Kind:        "EKS",
				Summary:     "Provides a production cluster within EKS",
				Description: "EKS Production Cluster",
				Labels: map[string]string{
					Label("environment"): "production",
					Label("kind"):        "EKS",
					Label("plural"):      "ekss",
				},
				Configuration: apiextv1.JSON{
					Raw: []byte(`{
						"authProxyAllowedIPs": [
							"0.0.0.0/0"
						],
						"defaultTeamRole": "view",
						"description": "eks-production cluster",
						"domain": "default",
						"enableDefaultTrafficBlock": false,
						"inheritTeamMembers": true,
						"privateIPV4Cidr": "10.0.0.0/16",
						"region": "eu-west-2",
						"version": "1.15",
						"nodeGroups": [
							{
								"name": "default",
								"instanceType": "c4.xlarge",
								"diskSize": 10,
								"enableAutoscaler": true,
								"minSize": 3,
								"desiredSize": 3,
								"maxSize": 12
							}
						]
					}`),
				},
			},
		},
	}
}
