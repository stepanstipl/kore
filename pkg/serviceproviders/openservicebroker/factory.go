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

package openservicebroker

import (
	"fmt"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
)

func init() {
	kore.RegisterServiceProviderFactory(ProviderFactory{})
}

type ProviderFactory struct{}

func (d ProviderFactory) Type() string {
	return "osb"
}

func (d ProviderFactory) JSONSchema() string {
	return `{
		"$id": "https://appvia.io/schemas/serviceprovider/osb.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "Open Service Broker Provider configuration schema",
		"type": "object",
		"additionalProperties": false,
		"required": [
			"url"
		],
		"properties": {
			"url": {
				"type": "string",
				"minLength": 1
			},
			"api_version": {
				"type": "string",
				"minLength": 1
			},
			"insecure": {
				"type": "boolean"
			},
			"ca_data": {
				"type": "string"
			},
			"auth_config": {
				"type": "object",
				"additionalProperties": false,
				"properties": {
					"basic_auth_config": {
						"type": "object",
						"additionalProperties": false,
						"required": [
							"username",
							"password"
						],
						"properties": {
							"username": {
								"type": "string",
								"minLength": 1
							},
							"password": {
								"type": "string"
							}
						}
					},
					"bearer_config": {
						"type": "object",
						"additionalProperties": false,
						"required": [
							"token"
						],
						"properties": {
							"token": {
								"type": "string",
								"minLength": 1
							}
						}
					}
				}
			}
		}
	}`
}

func (d ProviderFactory) CreateProvider(serviceProvider servicesv1.ServiceProvider) (kore.ServiceProvider, error) {
	var config = osb.DefaultClientConfiguration()
	config.Name = serviceProvider.Name

	if err := serviceProvider.Spec.GetConfiguration(config); err != nil {
		return nil, fmt.Errorf("failed to process service provider configuration: %w", err)
	}

	osbClient, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	provider, err := NewProvider(serviceProvider.Name, osbClient)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
