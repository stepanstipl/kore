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

package application

const ProviderSchema = `
{
	"$id": "https://appvia.io/schemas/serviceprovider/application.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Kubernetes Application provider",
	"type": "object",
	"additionalProperties": false
}`

const AppSchema = `
{
    "$id": "https://appvia.io/schemas/services/application/application.json",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "description": "Kubernetes Application",
    "type": "object",
    "additionalProperties": false,
    "required": [
        "resources"
    ],
    "properties": {
        "resources": {
            "type": "string",
            "format": "multiline"
        },
        "values": {
            "type": "string",
            "format": "multiline"
        }
    }
}
`

const HelmAppSchema = `
{
	"$id": "https://appvia.io/schemas/services/application/application.json",
	"$schema": "http://json-schema.org/draft-07/schema#",
	"description": "Kubernetes Application",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"source"
	],
	"properties": {
		"source": {
			"type": "object",
			"additionalProperties": false,
			"minProperties": 1,
			"maxProperties": 1,
			"properties": {
				"git": {
					"type": "object",
					"required": [
						"url",
						"ref"
					],
					"additionalProperties": false,
					"properties": {
						"url": {
							"type": "string",
							"minLength": 1,
							"description": "The URL of the Git repository",
							"examples": ["http://github.com/org/repo"]
						},
						"ref": {
							"type": "string",
							"minLength": 1,
							"default": "master",
							"description": "The Git branch (or other reference) to use",
							"examples": ["master", "v1.2.3"]
						},
						"path": {
							"type": "string",
							"description": "The path to the chart relative to the repository root",
							"examples": ["charts/my_chart"]
						}
					}
				},
				"helm": {
					"type": "object",
					"required": [
						"url",
						"name",
						"version"
					],
					"additionalProperties": false,
					"properties": {
						"url": {
							"type": "string",
							"minLength": 1,
							"description": "The URL of the Helm repository",
							"examples": ["https://charts.example.com"]
						},
						"name": {
							"type": "string",
							"minLength": 1,
							"description": "The name of the Helm chart (without an alias)"
						},
						"version": {
							"type": "string",
							"minLength": 1,
							"description": "The targeted Helm chart version",
							"examples": ["1.2.3"]
						}
					}
				}
			}
		},
		"values": {
			"type": "string",
			"format": "multiline"
		},
		"resourceKinds": {
			"type": "array",
			"items": {
				"type": "object",
				"required": [
					"kind"
				],
				"additionalProperties": false,
				"properties": {
					"group": {
						"type": "string",
						"description": "Kubernetes API group",
						"examples": ["apps"]
					},
					"kind": {
						"type": "string",
						"minLength": 1,
						"description": "Kubernetes resource kind",
						"examples": ["Deployment", "Service"]
					}
				}
			}
		},
		"resourceSelector": {
			"type": "object",
			"required": [
				"matchLabels"
			],
			"additionalProperties": false,
			"properties": {
				"matchLabels": {
					"type": "object",
					"additionalProperties": {
						"type": "string"
					}
				}
			}
		}
	}
}
`
