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

package security

import (
	"context"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Scanner represents the top-level interface to the security rule scanner
type Scanner interface {
	// ScanPlan scans the given plan against the registered rules and returns the result
	// of that scan
	ScanPlan(ctx context.Context, client client.Client, target *configv1.Plan) *securityv1.SecurityScanResult

	// ScanCluster scans the given cluster against the registered rules and returns the result
	// of that scan
	ScanCluster(ctx context.Context, client client.Client, target *clustersv1.Cluster) *securityv1.SecurityScanResult

	// GetRules returns all rules registered with this scanner
	GetRules() []Rule

	// GetRule returns a specific rule by its unique code from the rule set registered
	// with this scanner.
	GetRule(code string) Rule

	// RegisterRule adds an additional rule to this scanner
	RegisterRule(rule Rule)
}

// Rule implementations can be executed against plans or clusters to check for
// compliance with a specific security concern
type Rule interface {
	// Code returns the unique code identifying this rule
	Code() string
	// Name returns the human-readable name of this rule
	Name() string
	// Description returns the markdown-formatted description of this rule
	Description() string
}

// PlanRule implementations can be executed against a plan
type PlanRule interface {
	// CheckPlan runs this rule against the specified plan
	CheckPlan(ctx context.Context, client client.Client, target *configv1.Plan) (*securityv1.SecurityScanRuleResult, error)
}

// ClusterRule implementations can be executed against a cluster
type ClusterRule interface {
	// CheckCluster runs this rule against the specified cluster
	CheckCluster(ctx context.Context, client client.Client, target *clustersv1.Cluster) (*securityv1.SecurityScanRuleResult, error)
}
