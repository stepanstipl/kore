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
	"errors"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"

	"github.com/tidwall/gjson"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GKEAutoscaling determines whether the auth proxy IP range is suitably limited
type GKEAutoscaling struct{}

// Implement Rule

// Code returns the idenfier for this rule
func (p *GKEAutoscaling) Code() string {
	return "GKE-02"
}

// Name returns the name of this rule
func (p *GKEAutoscaling) Name() string {
	return "GKE Autoscaling"
}

// Description returns the markdown-formatted description of this rule
func (p *GKEAutoscaling) Description() string {
	return `
## Overview

This rule checks the status of the autoscaling on the GKE plans or clusters.

## Details

Autoscaling on GKE permits the control plan to scale the nodegroups based on the requirements. As the compute needs
(cpu, memory and scheduling constraints) increase on a GKE cluster, they can scale up the nodegroups up to meet the
requirement and or course back down when no longer required.

## Impact of warnings from this rule

Having the feature disabled means the cluster will have to manually scaled otherwise you might hit scheduling issues
due to a lack of nodes.
`
}

// ensureFeature handles the feature for both plans anc clusters
func (p *GKEAutoscaling) ensureFeature(config string) (*securityv1.SecurityScanRuleResult, error) {
	result := &securityv1.SecurityScanRuleResult{
		RuleCode: p.Code(),
		Status:   securityv1.Warning,
	}
	feature := gjson.Get(string(config), "enableAutoscaler")
	if !feature.Exists() {
		result.Message = "Could not check cluster due to invalid JSON"

		return result, errors.New("enableAutoscaler parameter does not exist")
	}

	// @step: check if the cluster version is out-of-date
	if !feature.Bool() {
		result.Message = "GKE Autoscaler is disabled "
		result.Status = securityv1.Warning

		return result, nil
	}

	result.Message = "GKE Autoscaling is enabled"
	result.Status = securityv1.Compliant

	return result, nil
}

// CheckPlan checks a plan for compliance with this rule
func (p *GKEAutoscaling) CheckPlan(ctx context.Context, cc client.Client, target *configv1.Plan) (*securityv1.SecurityScanRuleResult, error) {
	if target.Spec.Kind != "GKE" {
		return nil, nil
	}

	return p.ensureFeature(string(target.Spec.Configuration.Raw))
}

// CheckCluster checks a cluster for compliance with this rule
func (p *GKEAutoscaling) CheckCluster(ctx context.Context, cc client.Client, target *clustersv1.Cluster) (*securityv1.SecurityScanRuleResult, error) {
	if target.Spec.Kind != "GKE" {
		return nil, nil
	}

	return p.ensureFeature(string(target.Spec.Configuration.Raw))
}
