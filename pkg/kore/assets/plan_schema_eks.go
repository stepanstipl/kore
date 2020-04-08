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
		"region",
		"roleARN",
		"securityGroupIDs",
		"subnetIDs",
		"version"
	],
	"properties": {
		"authProxyAllowedIPs": {
			"type": "array",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"minItems": 1
		},
		"authProxyImage": {
			"type": "string"
		},
		"clusterUsers": {
			"type": "array",
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
			"type": "string"
		},
		"description": {
			"type": "string",
			"minLength": 1
		},
		"domain": {
			"type": "string",
			"minLength": 1
		},
		"enableDefaultTrafficBlock": {
			"type": "boolean"
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
					"eC2SSHKey",
					"instanceType",
					"maxSize",
					"minSize",
					"name",
					"nodeIAMRole",
					"region",
					"subnets"
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
					"nodeIAMRole": {
						"type": "string",
						"minLength": 1
					},
					"region": {
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
					"subnets": {
						"type": "array",
						"items": {
							"type": "string",
							"minLength": 1
						},
						"minItems": 1
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
		"region": {
			"type": "string",
			"minLength": 1
		},
		"roleARN": {
			"type": "string",
			"minLength": 1
		},
		"securityGroupIDs": {
			"type": "array",
			"items": {
				"type": "string",
				"minLength": 1
			},
			"minItems": 1
		},
		"subnetIDs": {
			"type": "array",
			"items": {
				"type": "string",
				"minLength": 1
			},
			"minItems": 1
		},
		"version": {
			"type": "string",
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
