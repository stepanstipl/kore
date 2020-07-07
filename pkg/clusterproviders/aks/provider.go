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

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
)

const Kind = "AKS"

func init() {
	kore.RegisterClusterProvider(Provider{})
}

type Provider struct {
}

func (p Provider) Type() string {
	return Kind
}

// PlanJSONSchema returns the JSON schema for the plans belonging to this cluster
func (p Provider) PlanJSONSchema() string {
	return schema
}

// DefaultPlans returns with the built-in default plans for the cluster provider
func (p Provider) DefaultPlans() []configv1.Plan {
	return plans
}

// DefaultPlanPolicy returns with the built-in default plan policy
func (p Provider) DefaultPlanPolicy() *configv1.PlanPolicy {
	return &planPolicy
}
