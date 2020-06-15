// +build !ignore_autogenerated

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1beta1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/appvia/kore/pkg/apis/monitoring/v1beta1.AlertSpec":   schema_pkg_apis_monitoring_v1beta1_AlertSpec(ref),
		"github.com/appvia/kore/pkg/apis/monitoring/v1beta1.AlertStatus": schema_pkg_apis_monitoring_v1beta1_AlertStatus(ref),
		"github.com/appvia/kore/pkg/apis/monitoring/v1beta1.RuleSpec":    schema_pkg_apis_monitoring_v1beta1_RuleSpec(ref),
	}
}

func schema_pkg_apis_monitoring_v1beta1_AlertSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AlertSpec specifies the details of a alert",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"event": {
						SchemaProps: spec.SchemaProps{
							Description: "Event is the raw event payload",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"summary": {
						SchemaProps: spec.SchemaProps{
							Description: "Summary is human readable summary for the alert",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"summary"},
			},
		},
	}
}

func schema_pkg_apis_monitoring_v1beta1_AlertStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AlertStatus is the status of the alert",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"archivedAt": {
						SchemaProps: spec.SchemaProps{
							Description: "ArchivedAt is indicates if the alert has been archived",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"detail": {
						SchemaProps: spec.SchemaProps{
							Description: "Detail provides a human readable message related to the current status of the alert",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"expiration": {
						SchemaProps: spec.SchemaProps{
							Description: "Expiration is the time the silience will finish",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the status of the alert",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"archivedAt"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}

func schema_pkg_apis_monitoring_v1beta1_RuleSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RuleSpec specifies the details of a alert rule",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"severity": {
						SchemaProps: spec.SchemaProps{
							Description: "Severity is the importance of the rule",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"source": {
						SchemaProps: spec.SchemaProps{
							Description: "Source is the provider of the rule i.e. prometheus, or a named source",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"summary": {
						SchemaProps: spec.SchemaProps{
							Description: "Summary is a summary of the rule",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"alerts": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Alerts is a collection of alerts related to the rule",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/monitoring/v1beta1.Alert"),
									},
								},
							},
						},
					},
					"rawRule": {
						SchemaProps: spec.SchemaProps{
							Description: "RawRule is the underlying rule definition",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"resource": {
						SchemaProps: spec.SchemaProps{
							Description: "Resource is the resource the alert is for",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
				},
				Required: []string{"severity", "source", "summary", "alerts", "rawRule", "resource"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Ownership", "github.com/appvia/kore/pkg/apis/monitoring/v1beta1.Alert"},
	}
}
