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

// EKSPlanSchema is the JSON schema used to describe and validate EKS Plans
const EKSPlanSchema = `
{
	"$id": "https://appvia.io/schemas/eks/plan.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "EKS Cluster Plan Schema",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"authProxyAllowedIPs",
		"description",
		"domain",
		"enableDefaultTrafficBlock",
		"inheritTeamMembers",
		"nodeGroups",
		"privateIPV4Cidr",
		"region",
		"version"
	],
	"properties": {
		"authorizedMasterNetworks": {
			"type": "array",
			"description": "A collection of network cidr allowed to speak the EKS control plan",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			}
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
		"domain": {
			"type": "string",
			"description": "The domain for this cluster.",
			"minLength": 1
		},
		"enableDefaultTrafficBlock": {
			"type": "boolean"
		},
		"fargateProfiles": {
			"type": "array",
			"description": "A collection of fargate profiles to provision",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"name",
					"subnets",
					"selectors",
				],
				"properties": {
					"name": {
						"description": "The name of the fargate profile.",
						"type": "string",
						"minLength": 1
					},
					"arn": {
						"description": "An optional pod execution IAM Arn the pods to run under.",
						"type": "string"
					},
					"selectors": {
						"type": "array",
						"description": "A collection of filers to match pods which should run on fargate.",
						"additionalProperties": false,
						"properties": {
							"namespace": {
								"description": "Selects all the pods within the namespace to run on fargate.",
								"type": "string",
								"minLength": 1
							},
							"labels": {
								"description": "Selects the pods based on matching labels.",
								"type": "object",
								"additionalProperties": { "type": "string" }
							}
						}
					}
				}
			}
		},
		"inheritTeamMembers": {
			"type": "boolean"
		},
		"nodeGroups": {
			"type": "array",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"desiredSize",
					"diskSize",
					"instanceType",
					"maxSize",
					"minSize",
					"name"
				],
				"properties": {
					"amiType": {
						"type": "string",
						"minLength": 1
					},
					"desiredSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1
					},
					"diskSize": {
						"type": "number",
						"multipleOf": 1
					},
					"eC2SSHKey": {
						"type": "string",
						"minLength": 1
					},
					"instanceType": {
						"type": "string",
						"minLength": 1
					},
					"labels": {
						"type": "object",
						"propertyNames": {
						  "minLength": 1,
						  "pattern": "^[a-zA-Z0-9\\-\\.\\_]+"
					    },
						"additionalProperties": { "type": "string" }
					},
					"minSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1
					},
					"maxSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1
					},
					"name": {
						"type": "string",
						"minLength": 1
					},
					"releaseVersion": {
						"type": "string",
						"minLength": 1
					},
					"sshSourceSecurityGroups": {
						"type": "array",
						"items": {
							"type": "string",
							"minLength": 1
						}
					},
					"tags": {
						"type": "object",
						"propertyNames": {
						  "minLength": 1,
						  "pattern": "^[a-zA-Z0-9+\\-=\\.\\_:/@]+"
					    },
						"additionalProperties": { "type": "string" }
					}
				}
			},
			"minItems": 1
		},
		"privateIPV4Cidr": {
			"type": "string",
			"description": "The range of IPv4 addresses for your EKS cluster in CIDR block format",
			"format": "1.2.3.4/16",
			"immutable": true
		},
		"region": {
			"type": "string",
			"description": "The AWS region in which this cluster will reside (e.g. eu-west-2).",
			"minLength": 1,
			"immutable": true
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
