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

import (
	"bytes"
	"encoding/json"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	"sigs.k8s.io/yaml"
)

// ServicePlan is the schema for the serviceplan
var ServicePlan = `
---
apiVersion: services.kore.appvia.io/v1
kind: ServicePlan
metadata:
  name: helm-app-kore-monitoring
  namespace: kore
  annotations:
    helm.values.schema: |
      {
        "$id": "https://appvia.io/schemas/kore-monitoring/plan.json",
        "$schema": "http://json-schema.org/draft-07/schema#",
        "description": "Kore Monitoring",
        "type": "object",
        "additionalProperties": false,
        "required": [
          "notifications"
        ],
        "properties": {
          "notifications": {
            "type": "object",
            "description": "Notification configuration is where to send the alerts",
            "additionalProperties": false,
            "properties": {
              "slack": {
                "type": "array",
                "descriptions": "Configuration for Slack notifications",
                "items": {
                  "type": "object",
                  "additionalProperties": false,
                  "required": [
                    "api_url"
                  ],
                  "properties": {
                    "api_url": {
                      "type": "string",
                      "description": "The slack wehbook url used to send notifications",
                      "minLength": 1
                    },
                    "channel": {
                      "type": "string",
                      "description": "An optional channel to send the message to",
                      "minLength": 1
                    },
                    "title": {
                      "type": "string",
                      "description": "The title to be used on the alerts"
                    },
                    "description": {
                      "type": "string",
                      "description": "An optional description for this notification channel"
                    },
                    "priority": {
                      "type": "array",
                      "description": "An optional priority for this notificaton",
                      "items": {
                        "type": "string",
                        "minLength": 1,
                        "enum": [
                          "P1",
                          "P2",
                          "P3",
                          "P4"
                        ]
                      },
                      "minItems": 1
                    },
                    "tags": {
                      "type": "array",
                      "description": "A collection of tags added to the notifications i.e. P1, P2",
                      "items": {
                        "type": "string",
                        "minLength": 1
                      }
                    }
                  }
                }
              },
              "pagerduty": {
                "type": "array",
                "description": "Configuration for notifications to Pagerduty",
                "items": {
                  "type": "object",
                  "additionalProperties": false,
                  "properties": {
                    "routing_key": {
                      "type": "string",
                      "description": "The PagerDuty integration key (when using PagerDuty integration type Events API v2)",
                      "minLength": 1
                    },
                    "service_key": {
                      "type": "string",
                      "description": "The PagerDuty integration key (when using PagerDuty integration type Prometheus)",
                      "minLength": 1
                    }
                  }
                }
              }
            }
          }
        }
      }
spec:
  kind: helm-app
  summary: Managaed Cluster Monitoring
  displayName: Kore Cluster Monitoring
  description: |
    Kore Monitoring service provides a monitoring stack used to ensure the health of the clusters and the services
  configuration:
    source:
      helm:
        url: https://storage.googleapis.com/kore-charts
        name: kore-monitoring
        version: 0.0.1
    resourceKinds:
      - group: apps
        kind: Deploymnent
      - group: ""
        kind: Service
`

// GetKoreMonitoringServicePlan returns the json encoded
func GetKoreMonitoringServicePlan() (*servicesv1.ServicePlan, error) {
	v, err := yaml.YAMLToJSON([]byte(ServicePlan))
	if err != nil {
		return nil, err
	}

	o := &servicesv1.ServicePlan{}
	if err := json.NewDecoder(bytes.NewReader(v)).Decode(o); err != nil {
		return nil, err
	}

	return o, nil
}
