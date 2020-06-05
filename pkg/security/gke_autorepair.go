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

// GKEAutorepair determines whether the auth proxy IP range is suitably limited
type GKEAutorepair struct{}

// Implement Rule

// Code returns the idenfier for this rule
func (p *GKEAutorepair) Code() string {
	return "GKE-03"
}

// Name returns the name of this rule
func (p *GKEAutorepair) Name() string {
	return "GKE Autorepair"
}

// Description returns the markdown-formatted description of this rule
func (p *GKEAutorepair) Description() string {
	return `
## Overview

This rule checks the status of the auto repair on the GKE plans or clusters.

## Details

Autorepair on GKE permits the control plan to automatically replace nodes which have failed the
kubernetes health checks.

## Impact of warnings from this rule

Not having this enabled means you may have unschedulable nodes or nodes with health concerned.
`
}

// ensureFeature handles the feature for both plans anc clusters
func (p *GKEAutorepair) ensureFeature(config string) (*securityv1.SecurityScanRuleResult, error) {
	// @TODO: Check all node pools where there are more than one. For now, just checking default node pool.
	return ValueAsExpected(p.Code(), config, "nodePools.0.enableAutorepair", true, securityv1.Warning,
		"GKE Autorepair is enabled",
		"GKE Autorepair is disabled",
	)
}

// CheckPlan checks a plan for compliance with this rule
func (p *GKEAutorepair) CheckPlan(ctx context.Context, cc client.Client, target *configv1.Plan) (*securityv1.SecurityScanRuleResult, error) {
	if target.Spec.Kind != "GKE" {
		return nil, nil
	}

	return p.ensureFeature(string(target.Spec.Configuration.Raw))
}

// CheckCluster checks a cluster for compliance with this rule
func (p *GKEAutorepair) CheckCluster(ctx context.Context, cc client.Client, target *clustersv1.Cluster) (*securityv1.SecurityScanRuleResult, error) {
	if target.Spec.Kind != "GKE" {
		return nil, nil
	}

	return p.ensureFeature(string(target.Spec.Configuration.Raw))
}
