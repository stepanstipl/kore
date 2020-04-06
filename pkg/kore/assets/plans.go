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
				Name: "gke-development",
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
				Values: apiextv1.JSON{
					Raw: []byte(`
						{
							"authorizedMasterNetworks": [
								{
									"name": "default",
									"cidr": "0.0.0.0/0"
								}
                            ],
							"authProxyAllowedIPs":           ["0.0.0.0/0"],
							"description":                   "gke-development cluster",
							"diskSize":                      100,
							"enableAutoupgrade":             true,
							"enableAutorepair":              true,
							"enableAutoscaler":              true,
							"enableHTTPLoadBalancer":        true,
							"enableHorizontalPodAutoscaler": true,
							"enableIstio":                   false,
							"enablePrivateEndpoint":         false,
							"enablePrivateNetwork":          false,
							"enableShieldedNodes":           true,
							"enableStackDriverLogging":      true,
							"enableStackDriverMetrics":      true,
							"imageType":                     "COS",
							"machineType":                   "n1-standard-2",
							"maintenanceWindow":             "03:00",
							"maxSize":                       10,
							"network":                       "default",
							"region":                        "europe-west2",
							"size":                          1,
							"subnetwork":                    "default",
							"version":                       "1.15.11-gke.1"
						}`,
					),
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gke-production",
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
				Values: apiextv1.JSON{
					Raw: []byte(`
						{
							"authorizedMasterNetworks": [
								{
									"name": "default",
									"cidr": "0.0.0.0/0"
								}
                            ],
							"authProxyAllowedIPs":           ["0.0.0.0/0"],
							"description":                   "gke-production cluster",
							"diskSize":                      100,
							"enableAutoupgrade":             true,
							"enableAutorepair":              true,
							"enableAutoscaler":              true,
							"enableHTTPLoadBalancer":        true,
							"enableHorizontalPodAutoscaler": true,
							"enableIstio":                   false,
							"enablePrivateEndpoint":         false,
							"enablePrivateNetwork":          false,
							"enableShieldedNodes":           true,
							"enableStackDriverLogging":      true,
							"enableStackDriverMetrics":      true,
							"imageType":                     "COS",
							"machineType":                   "n1-standard-2",
							"maintenanceWindow":             "03:00",
							"maxSize":                       10,
							"network":                       "default",
							"region":                        "europe-west2",
							"size":                          2,
							"subnetwork":                    "default",
							"version":                       "1.15.11-gke.1"
						}`,
					),
				},
			},
		},
	}
}
