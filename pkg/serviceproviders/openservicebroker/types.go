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
	osb "github.com/kubernetes-sigs/go-open-service-broker-client/v2"
)

type ProviderConfiguration struct {
	osb.ClientConfiguration `json:",inline"`
	// DefaultPlanNames is a list of plans to use as default for each service kind in a format as `[kind]-[plan name]`
	DefaultPlanNames []string `json:"defaultPlanNames,omitempty"`
	// AllowEmptyCredentialSchema will allow plans with empty credential schemas
	AllowEmptyCredentialSchema bool `json:"allowEmptyCredentialSchema"`
}

// ProviderData will store the "operation" value returned from async operations
type ProviderData struct {
	Operation *osb.OperationKey `json:"operation,omitempty"`
}

type ServiceKindProviderData struct {
	// ServiceID is the service kind identifier in the service provider
	ServiceID string `json:"serviceID,omitempty"`
	// DefaultPlanID is the default plan id which is used for user-created service plans
	DefaultPlanID string `json:"defaultPlanID,omitempty"`
}

type ServicePlanProviderData struct {
	// PlanID is the service plan identifier in the service provider
	PlanID string `json:"planID,omitempty"`
	// ServiceID is the service kind identifier in the service provider
	ServiceID string `json:"serviceID,omitempty"`
}
