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

package aks

//go:generate go run github.com/appvia/kore/cmd/struct-gen Configuration
const schema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"authorizedMasterNetworks",
		"authProxyAllowedIPs",
		"description",
		"dnsPrefix",
		"domain",
		"networkPlugin",
		"nodePools",
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
			"minLength": 1,
			"description": "Geographical location for this cluster",
			"examples": ["ukwest", "eastus"],
			"immutable": true
		},
		"version": {
			"type": "string",
			"description": "Kubernetes version",
			"pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$",
			"examples": [
				"1.16.10", "1.17.7"
			],
			"default": ""
		},
		"dnsPrefix": {
			"type": "string",
			"description": "DNS name prefix to use with the hosted Kubernetes API server FQDN.",
			"minLength": 1,
			"immutable": true
		},
		"nodePools": {
			"type": "array",
			"minItems": 1,
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"name",
					"mode",
					"size",
					"machineType",
					"imageType",
					"diskSize"
				],
				"properties": {
					"name": {
						"type": "string",
						"pattern": "^[a-z][-a-z0-9]{0,38}[a-z0-9]$",
						"description": "Name of this node pool. Must be unique within the cluster.",
						"immutable": true
					},
					"mode": {
						"type": "string",
						"enum": ["System", "User"],
						"description": "Type of the node pool.\nSystem node pools serve the primary purpose of hosting critical system pods such as CoreDNS and tunnelfront.\nUser node pools serve the primary purpose of hosting your application pods."
					},
					"version": {
						"type": "string",
						"description": "Node pool version, blank to use same version as cluster (recommended).",
						"default": ""
					},
					"enableAutoscaler": {
						"type": "boolean",
						"default": true,
						"description": "Add and remove nodes automatically based on load"
					},
					"minSize": {
						"type": "integer",
						"minimum": 1,
						"maximum": 100,
						"default": 1,
						"description": "The minimum nodes this pool should contain (if auto-scale enabled)"
					},
					"maxSize": {
						"type": "integer",
						"minimum": 1,
						"maximum": 100,
						"default": 10,
						"description": "The maximum nodes this pool should contain (if auto-scale enabled)"
					},
					"size": {
						"type": "integer",
						"minimum": 1,
						"maximum": 100,
						"default": 1,
						"description": "How many nodes to build when provisioning this pool - if autoscaling enabled, this will be the initial size",
						"immutable": true
					},
					"maxPodsPerNode": {
						"type": "number",
						"multipleOf": 1,
						"description": "The maximum number of pods that can be scheduled onto each node of this pool - if left blank, it will set this automatically based on the machine type",
						"maximum": 110,
						"minimum": 8,
						"immutable": true
					},
					"machineType": {
						"type": "string",
						"minLength": 1,
						"description": "The type of nodes used for this node pool",
						"default": "Standard_D1_v2",
						"immutable": true
					},
					"imageType": {
						"type": "string",
						"enum": [ "Linux", "Windows" ],
						"description": "The image type used by the nodes",
						"default": "Linux"
					},
					"diskSize": {
						"type": "number",
						"description": "The amount of storage in GiB provisioned on the nodes in this group",
						"multipleOf": 1,
						"default": 100,
						"minimum": 30,
						"maximum": 65536,
						"immutable": true
					},
					"labels": {
						"type": "object",
						"propertyNames": {
						  "minLength": 1,
						  "pattern": "^[a-zA-Z0-9\\-\\.\\_]+$"
					    },
						"additionalProperties": { "type": "string" },
						"description": "A set of labels to help Kubernetes workloads find this group",
						"default": {},
						"immutable": true
					},
					"taints": {
						"type": "array",
						"description": "A collection of kubernetes taints to add on the nodes.",
						"default": [],
						"immutable": true,
						"items": {
							"type": "object",
							"additionalProperties": false,
							"required": [
								"effect",
								"key",
								"value"
							],
							"properties": {
								"key": {
									"type": "string",
									"minLength": 1,
									"description": "Taint key",
									"immutable": true
								},
								"value": {
									"type": "string",
									"minLength": 1,
									"description": "Taint value",
									"immutable": true
								},
								"effect": {
									"type": "string",
									"enum": [ "NoSchedule", "PreferNoSchedule", "NoExecute", "NoEffect" ],
									"description": "The chosen effect of the taint",
									"default": "NoSchedule",
									"immutable": true
								}
							}
						}
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
			}
		},
		"defaultTeamRole": {
			"type": "string",
			"description": "The default role that team members have on this cluster.",
			"enum": [ "view", "edit", "admin", "cluster-admin" ],
			"default": "view"
		},
		"inheritTeamMembers": {
			"type": "boolean"
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
		"domain": {
			"type": "string",
			"description": "The domain for this cluster.",
			"minLength": 1,
			"default": "default"
		},
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
			"minItems": 1,
			"default": [
				{ "name": "default", "cidr": "0.0.0.0/0" }
			]
		},
		"authProxyAllowedIPs": {
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
			"type": "boolean"
		},
		"enablePodSecurityPolicy": {
			"type": "boolean"
		},
		"networkPlugin": {
			"type": "string",
			"enum": [
				"azure",
				"kubenet"
			],
			"immutable": true
		},
		"networkPolicy": {
			"type": "string",
			"enum": [
				"",
				"azure",
				"calico"
			],
			"immutable": true
		},
		"privateClusterEnabled": {
			"type": "boolean",
			"immutable": true
		},
		"linuxProfile": {
			"type": "object",
			"required": [
				"adminUsername",
				"sshPublicKeys"
			],
			"additionalProperties": false,
			"properties": {
				"adminUsername": {
					"type": "string"
				},
				"sshPublicKeys": {
					"items": {
						"type": "string"
					},
					"type": "array"
				}
			},
			"immutable": true
		},
		"windowsProfile": {
			"type": "object",
			"required": [
				"adminUsername",
				"adminPassword"
			],
			"additionalProperties": false,
			"properties": {
				"adminPassword": {
					"type": "string"
				},
				"adminUsername": {
					"type": "string"
				}
			}
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
