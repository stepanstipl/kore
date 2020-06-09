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
			"title": "Auth Proxy Allowed IP Ranges",
			"type": "array",
			"description": "The networks which are allowed to connect to this cluster (e.g. via kubectl).",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"minItems": 1
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
							"minLength": 1,
							"enum": [ "view", "edit", "admin", "cluster-admin" ]
						},
						"minItems": 1
					}
				}
			}
		},
		"defaultTeamRole": {
			"type": "string",
			"description": "The default role that team members have on this cluster.",
			"enum": [ "view", "edit", "admin", "cluster-admin" ]
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
						"title": "Compute Type",
						"description": "Whether this node group is for general purpose or GPU workloads",
						"type": "string",
						"enum": ["AL2_x86_64", "AL2_x86_64_GPU"],
						"default": "AL2_x86_64"
					},
					"desiredSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 1
					},
					"diskSize": {
						"type": "number",
						"multipleOf": 1,
						"default": 10
					},
					"eC2SSHKey": {
						"title": "EC2 SSH Key",
						"description": "Reference to an key which exists in your AWS account to allow SSH access to nodes",
						"type": "string",
						"minLength": 1
					},
					"instanceType": {
						"type": "string",
						"minLength": 1,
						"default": "t3.medium"
					},
					"labels": {
						"type": "object",
						"propertyNames": {
						  "pattern": "^[a-zA-Z0-9\\-\\.\\_]+$"
					    },
						"additionalProperties": { "type": "string" },
						"description": "A set of labels to help Kubernetes workloads find this group",
						"default": {}
					},
					"minSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 1
					},
					"maxSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 10
					},
					"name": {
						"type": "string",
						"minLength": 1
					},
					"releaseVersion": {
						"type": "string",
						"description": "Blank to use latest (recommended), if set must be for same Kubernetes version as the top-level plan version and for the same AMI type as specified for this node group.",
						"pattern": "^($|[0-9]+\\.[0-9]+\\.[0-9]+\\-[0-9]+)$",
						"examples": [
							"1.16.8-20200507", "1.15.11-20200507"
						],
						"default": ""
					},
					"sshSourceSecurityGroups": {
						"title": "SSH Security Groups",
						"description": "Reference to security groups from which SSH access is permitted - must exist and be in same region as this cluster",
						"type": "array",
						"items": {
							"type": "string",
							"pattern": "^([0-9]*\\/)?sg-[0-9]+$"
						},
						"examples": [
							"sg-0123456789 (security group in same account as cluster)",
							"12345/sg-012346789 (security group in account 12345)"
						],
						"default": []
					},
					"tags": {
						"type": "object",
						"propertyNames": {
						  "pattern": "^[a-zA-Z0-9+\\-=\\.\\_:/@]+$"
					    },
						"additionalProperties": { "type": "string" },
						"default": {}
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
			"description": "The AWS region in which this cluster will reside",
			"examples": [ "eu-west-2", "us-east-1" ],
			"pattern": "^(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\\d$",
			"immutable": true
		},
		"version": {
			"type": "string",
			"description": "The Kubernetes version to deploy.",
			"pattern": "^[0-9]+\\.[0-9]+$",
			"examples": [
				"1.15", "1.16"
			]
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
		"required": ["defaultTeamRole"]
	},
	"else": {
	}
}
`
