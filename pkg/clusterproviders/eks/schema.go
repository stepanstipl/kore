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

//go:generate go run github.com/appvia/kore/cmd/struct-gen Configuration
const schema = `
{
	"$id": "https://appvia.io/kore/schemas/eks/plan.json",
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
		"description": {
			"type": "string",
			"description": "Meaningful description of this cluster.",
			"minLength": 1
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
						"default": 1,
						"description": "The minimum nodes this group should contain (if auto-scale enabled)"
					},
					"maxSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 10,
						"description": "The maximum nodes this group should contain (if auto-scale enabled)"
					},
					"enableAutoscaler": {
						"type": "boolean",
						"default": true,
						"description": "Will enable the cluster autoscaler to scale this specific nodegroup"
					},
					"name": {
						"type": "string",
						"minLength": 1,
						"immutable": true,
						"identifier": true
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
				},
				"allOf": [
					{
						"$comment": "Require min/max sizes if auto-scale enabled",
						"if": {
							"properties": {
								"enableAutoscaler": {
									"const": true
								}
							}
						},
						"then": {
							"properties": {
								"minSize": {
									"minimum": 1
								},
								"maxSize": {
									"minimum": 1
								}
							},
							"required": ["minSize", "maxSize"]
						}
					}
				]
			},
			"minItems": 1
		},
		"privateIPV4Cidr": {
			"type": "string",
			"description": "The range of IPv4 addresses for your EKS cluster in CIDR block format",
			"format": "1.2.3.4/16",
			"immutable": true,
			"default": "10.0.0.0/16"
		},
		"defaultTeamRole": {
			"type": "string",
			"description": "The role that team members will have on this cluster if 'inherit team members' enabled",
			"enum": [ "view", "edit", "admin", "cluster-admin" ],
			"default": "view"
		},
		"inheritTeamMembers": {
			"type": "boolean",
			"description": "Whether team members will all have access to this cluster by default",
			"default": true
		},
		"clusterUsers": {
			"type": "array",
			"description": "Users who should be allowed to access this cluster, will override any default role set above for these users",
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
		"domain": {
			"type": "string",
			"description": "The domain for this cluster.",
			"minLength": 1,
			"default": "default"
		},
		"authorizedMasterNetworks": {
			"type": "array",
			"description": "A collection of network cidr allowed to speak the EKS control plan",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"default": [ "0.0.0.0/0" ]
		},
		"authProxyAllowedIPs": {
			"title": "Auth Proxy Allowed IP Ranges",
			"type": "array",
			"description": "The networks which are allowed to connect to this cluster (e.g. via kubectl).",
			"items": {
				"type": "string",
				"format": "1.2.3.4/16"
			},
			"minItems": 1,
			"default": [ "0.0.0.0/0" ]
		},
		"enableDefaultTrafficBlock": {
			"type": "boolean",
			"default": false
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
