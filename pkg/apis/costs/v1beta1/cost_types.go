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

package v1beta1

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CostList represents a collection of costs
// +k8s:openapi-gen=false
type CostList struct {
	Items []Cost `json:"items"`
}

// Cost defines the details about the cost for a piece of infrastructure deployed by Kore for
// a team
// +k8s:openapi-gen=false
type Cost struct {
	// Resource is a reference to the piece of team infrastructure which this cost applies to.
	Resource corev1.Ownership `json:"resource,omitempty"`
	// Team is the name of the team this cost applies to.
	Team string `json:"team,omitempty"`
	// Cost is the actual incurred cost total cost for this piece of infrastructure for the
	// specified time period in microdollars
	Cost int64 `json:"cost,omitempty"`
	// From indicates the start of the period this cost is applicable for
	From metav1.Time `json:"from,omitempty"`
	// To indicates the end of the period this cost is applicable for
	To metav1.Time `json:"to,omitempty"`
	// RetrievedAt indicates the time this cost was retrieved from the provider by Kore
	RetrievedAt metav1.Time `json:"preparedAt,omitempty"`
	// CostElements provides details of the different components which make up this cost,
	// may be empty if the top level infrastructure does not have any sub-components.
	CostElements []CostElement `json:"costElements,omitempty"`
}

// CostElement represents a logical component which has an associated cost
// +k8s:openapi-gen=false
type CostElement struct {
	// Name is the name of this component
	Name string `json:"name,omitempty"`
	// Cost is the actual incurred cost in microdollars
	Cost int64 `json:"cost,omitempty"`
}
