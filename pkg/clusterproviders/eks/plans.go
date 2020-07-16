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

package eks

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var plans = []configv1.Plan{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "eks-development",
			Annotations: map[string]string{
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Provides a development cluster within EKS",
			Description: "EKS Development Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "development",
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
			Kind:        Kind,
			Summary:     "Provides a production cluster within EKS",
			Description: "EKS Production Cluster",
			Labels: map[string]string{
				kore.Label("environment"): "production",
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
