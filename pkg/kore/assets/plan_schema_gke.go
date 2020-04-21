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

// GKEPlanSchema is the JSON schema used to describe and validate GKE Plans
const GKEPlanSchema = `
{
	"$id": "https://appvia.io/schemas/gke/plan.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "GKE Cluster Plan Schema",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"authorizedMasterNetworks",
		"authProxyAllowedIPs",
		"description",
		"diskSize",
		"domain",
		"enableAutoupgrade",
		"enableAutorepair",
		"enableAutoscaler",
		"enableDefaultTrafficBlock",
		"enableHTTPLoadBalancer",
		"enableHorizontalPodAutoscaler",
		"enableIstio",
		"enablePrivateEndpoint",
		"enablePrivateNetwork",
		"enableShieldedNodes",
		"enableStackDriverLogging",
		"enableStackDriverMetrics",
		"imageType",
		"inheritTeamMembers",
		"machineType",
		"maintenanceWindow",
		"maxSize",
		"network",
		"region",
		"size",
		"subnetwork",
		"version"
	],
	"properties": {
		"authorizedMasterNetworks": {
			"type": "array",
			"description": "The networks which are allowed to access the master control plane.",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"name",
					"cidr"
				],
				"properties": {
					"name": {
						"type": "string",
						"minLength": 1
					},
					"cidr": {
						"type": "string",
						"format": "1.2.3.4/16"
					}
				}
			},
			"minItems": 1
		},
		"authProxyAllowedIPs": {
			"type": "array",
			"description": "The networks which are allowed to connect to this cluster (e.g. via kubectl).",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"minItems": 1
		},
		"authProxyImage": {
			"type": "string",
			"description": "TBC"
		},
		"clusterUsers": {
			"type": "array",
			"description": "Users who should be allowed to access this cluster.",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"username",
					"roles"
				],
				"properties": {
					"username": {
						"type": "string",
						"minLength": 1
					},
					"roles": {
						"type": "array",
						"items": {
							"type": "string",
							"minLength": 1
						},
						"minItems": 1
					}
				}
			}
		},
		"defaultTeamRole": {
			"type": "string",
			"description": "The default role that team members have on this cluster."
		},
		"description": {
			"type": "string",
			"description": "Meaningful description of this cluster.",
			"minLength": 1
		},
		"diskSize": {
			"type": "number",
			"description": "The amount of storage provisioned on this cluster? Or on its nodes?",
			"multipleOf": 1,
			"minimum": 10,
			"maximum": 65536
		},
		"domain": {
			"type": "string",
			"description": "The domain for this cluster.",
			"minLength": 1
		},
		"enableAutoupgrade": {
			"type": "boolean",
			"description": "Enable to keep this cluster updated whenever new versions of Kubernetes are made available by GCP."
		},
		"enableAutorepair": {
			"type": "boolean"
		},
		"enableAutoscaler": {
			"type": "boolean"
		},
		"enableDefaultTrafficBlock": {
			"type": "boolean"
		},
		"enableHTTPLoadBalancer": {
			"type": "boolean"
		},
		"enableHorizontalPodAutoscaler": {
			"type": "boolean"
		},
		"enableIstio": {
			"type": "boolean"
		},
		"enablePrivateEndpoint": {
			"type": "boolean"
		},
		"enablePrivateNetwork": {
			"type": "boolean"
		},
		"enableShieldedNodes": {
			"type": "boolean"
		},
		"enableStackDriverLogging": {
			"type": "boolean"
		},
		"enableStackDriverMetrics": {
			"type": "boolean"
		},
		"imageType": {
			"type": "string",
			"minLength": 1
		},
		"inheritTeamMembers": {
			"type": "boolean"
		},
		"machineType": {
			"type": "string",
			"description": "The type of nodes used for this cluster's node pool.",
			"minLength": 1
		},
		"maintenanceWindow": {
			"type": "string",
			"description": "Time of day to allow maintenance operations to be performed by the cloud provider on this cluster.",
			"format": "hh:mm"
		},
		"maxSize": {
			"type": "number",
			"description": "Maximum number of nodes? to allow for this cluster if auto-scaling enabled.",
			"multipleOf": 1,
			"minimum": 0
		},
		"network": {
			"type": "string",
			"minLength": 1
		},
		"region": {
			"type": "string",
			"minLength": 1
		},
		"size": {
			"type": "number",
			"multipleOf": 1,
			"minimum": 0
		},
		"subnetwork": {
			"type": "string",
			"minLength": 1
		},
		"version": {
			"type": "string",
			"description": "The Kubernetes version to deploy.",
			"minLength": 1
		}
	},
	"if": {
		"properties": {
			"inheritTeamMembers": {
				"const": true
			}
		},
		"required": ["inheritTeamMembers"]
	},
	"then": {
		"properties": {
			"defaultTeamRole": {
				"minLength": 1
			}
		},
		"required": ["defaultTeamRole"]
	},
	"else": {
	}
}
`
