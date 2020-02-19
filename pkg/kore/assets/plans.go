/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
							"diskSize":               100,
							"enableAutorepair":       true,
							"enableAutoscaler":       true,
							"enableHTTPLoadBalancer": true,
							"enableHorizontalPodAutoscaler": true,
							"enableIstio": false,
							"enablePrivateNetwork":   true,
							"enableStackDriverLogging": true,
							"enableStackDriverMetrics": true,
							"imageType":              "COS",
							"machineType":            "n1-standard-2",
							"maintenanceWindow":      "03:00",
							"maxSize":                10,
							"network":                "default",
							"size":                   1,
							"subnetwork":             "default",
							"version":                "1.14.9-gke.23"
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
							"diskSize":               100,
							"enableAutorepair":       true,
							"enableAutoscaler":       true,
							"enableHTTPLoadBalancer": true,
							"enableHorizontalPodAutoscaler": true,
							"enableIstio": false,
							"enablePrivateNetwork":   true,
							"enableStackDriverLogging": true,
							"enableStackDriverMetrics": true,
							"imageType":              "COS",
							"machineType":            "n1-standard-2",
							"maintenanceWindow":      "03:00",
							"masterIPV4Cidr": "172.16.0.0/28",
							"maxSize":                10,
							"network":                "default",
							"size":                   2,
							"subnetwork":             "default",
							"version":                "1.14.9-gke.23"
						}`,
					),
				},
			},
		},
	}
}
