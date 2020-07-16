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

package gke

const schema = `
{
	"$id": "https://appvia.io/kore/schemas/gke/plan.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "GKE Cluster Plan Schema",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"authorizedMasterNetworks",
		"authProxyAllowedIPs",
		"description",
		"domain",
		"enableDefaultTrafficBlock",
		"enableHTTPLoadBalancer",
		"enableHorizontalPodAutoscaler",
		"enableIstio",
		"enablePrivateEndpoint",
		"enablePrivateNetwork",
		"enableShieldedNodes",
		"enableStackDriverLogging",
		"enableStackDriverMetrics",
		"inheritTeamMembers",
		"maintenanceWindow",
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
			"examples": ["europe-west2", "us-east1"],
			"immutable": true
		},
		"releaseChannel": {
			"type": "string",
			"description": "Follow a GKE release channel to control the auto-upgrade of your cluster - if set, auto-upgrade will be true on all node groups",
			"enum": ["REGULAR", "STABLE", "RAPID", ""],
			"default": "REGULAR"
		},
		"version": {
			"type": "string",
			"description": "Kubernetes version - must be blank if release channel specified.",
			"pattern": "^($|-|latest|[0-9]+\\.[0-9]+($|\\.[0-9]+($|\\-gke\\.[0-9]+)))$",
			"examples": [
				"- (GKE default)", "1.15 (latest 1.15.x)", "1.15.1", "1.15.1-gke.6 (exact GKE patch version, not recommended)", "latest"
			],
			"default": ""
		},
		"nodePools": {
			"type": "array",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"name",
					"enableAutoscaler",
					"enableAutoupgrade",
					"enableAutorepair",
					"version",
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
					"enableAutoupgrade": {
						"type": "boolean",
						"description": "Enable to update this node pool updated when new GKE versions are made available by GCP - must be enabled if a release channel is selected",
						"default": true
					},
					"version": {
						"type": "string",
						"description": "Node pool version, blank to use same version as cluster (recommended); must be blank if cluster follows a release channel. Must be within 2 minor versions of the master version (e.g. for master version 1.16, this must be 1.14, 1.15 or 1.16) or 1 minor version if auto-upgrade enabled",
						"pattern": "^($|latest|[0-9]+\\.[0-9]+($|\\.[0-9]+($|\\-gke\\.[0-9]+)))$",
						"default": "",
						"examples": [
							"latest", "1.15 (latest 1.15.x-gke.y)", "1.15.1 (latest 1.15.1-gke.x)", "1.15.1-gke.6 (exact GKE version)"
						]
					},
					"enableAutoscaler": {
						"type": "boolean",
						"default": true,
						"description": "Add and remove nodes automatically based on load"
					},
					"enableAutorepair": {
						"type": "boolean",
						"default": true,
						"description": "Automatically repair any failed nodes within this node pool."
					},
					"minSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 1,
						"description": "The minimum nodes this pool should contain (if auto-scale enabled)"
					},
					"maxSize": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 10,
						"description": "The maximum nodes this pool should contain (if auto-scale enabled)"
					},
					"size": {
						"type": "number",
						"multipleOf": 1,
						"minimum": 1,
						"default": 1,
						"description": "How many nodes to build when provisioning this pool - if autoscaling enabled, this will be the initial size",
						"immutable": true
					},
					"maxPodsPerNode": {
						"type": "number",
						"multipleOf": 1,
						"description": "The maximum number of pods that can be scheduled onto each node of this pool",
						"default": 110,
						"maximum": 110,
						"minimum": 8,
						"immutable": true
					},
					"machineType": {
						"type": "string",
						"description": "The type of nodes used for this node pool",
						"pattern": "^[a-z][0-9]\\-(micro|small|medium|standard\\-[0-9]+|highmem\\-[0-9]+|highcpu\\-[0-9]+|ultramem\\-[0-9]+|megamem\\-[0-9]+)$",
						"default": "n1-standard-2",
						"immutable": true
					},
					"imageType": {
						"type": "string",
						"enum": [ "COS", "COS_CONTAINERD", "UBUNTU", "UBUNTU_CONTAINERD", "WINDOWS_LTSC", "WINDOWS_SAC" ],
						"description": "The image type used by the nodes",
						"default": "COS"
					},
					"diskSize": {
						"type": "number",
						"description": "The amount of storage in GiB provisioned on the nodes in this group",
						"multipleOf": 1,
						"default": 100,
						"minimum": 10,
						"maximum": 65536,
						"immutable": true
					},
					"preemptible": {
						"type": "boolean",
						"description": "Whether to use pre-emptible nodes (cheaper, but can and will be terminated at any time, use with care).",
						"default": false,
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
									"immutable": true
								},
								"value": {
									"type": "string",
									"minLength": 1,
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
		},
		"enablePrivateNetwork": {
			"type": "boolean",
			"immutable": true,
			"default": false
		},
		"enablePrivateEndpoint": {
			"type": "boolean",
			"immutable": true,
			"default": false
		},
		"enableHTTPLoadBalancer": {
			"type": "boolean",
			"default": true
		},
		"enableHorizontalPodAutoscaler": {
			"type": "boolean",
			"immutable": true,
			"default": true
		},
		"enableIstio": {
			"type": "boolean",
			"immutable": true,
			"default": false
		},
		"enableShieldedNodes": {
			"type": "boolean",
			"description": "Shielded nodes provide additional verifications of the node OS and VM, with enhanced rootkit and bootkit protection applied",
			"immutable": true,
			"default": true
		},
		"enableStackDriverLogging": {
			"type": "boolean",
			"immutable": true,
			"default": true
		},
		"enableStackDriverMetrics": {
			"type": "boolean",
			"immutable": true,
			"default": true
		},
		"maintenanceWindow": {
			"type": "string",
			"description": "Time of day to allow maintenance operations to be performed by the cloud provider on this cluster.",
			"format": "hh:mm",
			"immutable": true,
			"default": "03:00"
		},


		"diskSize": {
			"deprecated": true,
			"type": "number",
			"description": "DEPRECATED: Set disk size on node pool instead",
			"multipleOf": 1,
			"minimum": 10,
			"maximum": 65536
		},
		"enableAutoupgrade": {
			"deprecated": true,
			"description": "DEPRECATED: Set auto-upgrade on node pool instead",
			"type": "boolean"
		},
		"enableAutorepair": {
			"deprecated": true,
			"description": "DEPRECATED: Set auto-repair on node pool instead",
			"type": "boolean"
		},
		"enableAutoscaler": {
			"deprecated": true,
			"description": "DEPRECATED: Set auto-scale on node pool instead",
			"type": "boolean"
		},
		"imageType": {
			"deprecated": true,
			"description": "DEPRECATED: Set image type on node pool instead",
			"type": "string",
			"minLength": 1
		},
		"machineType": {
			"deprecated": true,
			"description": "DEPRECATED: Set machine type on node pool instead",
			"type": "string",
			"minLength": 1
		},
		"maxSize": {
			"deprecated": true,
			"description": "DEPRECATED: Set max size on node pool instead",
			"type": "number",
			"multipleOf": 1,
			"minimum": 0
		},
		"size": {
			"deprecated": true,
			"description": "DEPRECATED: Set size on node pool instead",
			"type": "number",
			"multipleOf": 1,
			"minimum": 0
		},
		"subnetwork": {
			"deprecated": true,
			"description": "DEPRECATED: Unused",
			"type": "string"
		},
		"network": {
			"deprecated": true,
			"description": "DEPRECATED: It is not supported to specify a custom network. This property will be ignored.",
			"type": "string",
			"minLength": 1,
			"immutable": true
		}
	},
	"allOf": [
		{
			"$comment": "Require default team role if inherit team members set",
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
			}
		},
		{
			"$comment": "If all deprecated fields not specified, make node pools and release channel required",
			"if": {
				"required": [
					"diskSize", "enableAutoupgrade", "enableAutorepair", "enableAutoscaler",
					"imageType", "machineType", "maxSize", "size", "version"
				]
			},
			"then": {
			},
			"else": {
				"properties": {
					"nodePools": {
						"minItems": 1
					}
				},
				"required": [ "nodePools", "releaseChannel" ]
			}
		},
		{
			"$comment": "Require auto-upgrade and no version on node pools if releaseChannel set, else require version if no release channel set",
			"if": {
				"properties": {
					"releaseChannel": {
						"const": ""
					}
				}
			},
			"then": {
				"properties": {
					"version": {
						"pattern": "^(-|latest|[0-9]+\\.[0-9]+($|\\.[0-9]+($|\\-gke\\.[0-9]+)))$"
					}
				}
			},
			"else": {
				"properties": {
					"version": {
						"const": ""
					},
					"nodePools": {
						"items": {
							"properties": {
								"enableAutoupgrade": {
									"const": true
								},
								"version": {
									"const": ""
								}
							}
						}
					}
				}
			}
		}
	]
}
`
