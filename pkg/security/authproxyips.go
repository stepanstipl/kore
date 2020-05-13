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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AuthProxyIPRangeRule determines whether the auth proxy IP range is suitably limited
type AuthProxyIPRangeRule struct {
}

// Implement Rule

// Code returns the idenfier for this rule
func (p *AuthProxyIPRangeRule) Code() string {
	return "AUTHIP-01"
}

// Name returns the name of this rule
func (p *AuthProxyIPRangeRule) Name() string {
	return "Auth Proxy IP Ranges"
}

// Description returns the markdown-formatted description of this rule
func (p *AuthProxyIPRangeRule) Description() string {
	return `
## Overview

This rule ensures that a plan specifies narrow (/16 or smaller) IP ranges for the authentication proxy to accept.

## Details

When Kore creates a Kubernetes cluster, it uses an authentication proxy running inside that cluster to authenticate
access to the cluster. It is best practice to restrict the IP address ranges enabled by default on a cluster to a
known set of IP ranges.

This rule returns a warning where a plan either does not specify any IP ranges, or specifies any IP range wider
than a /16.

## Impact of warnings from this rule

The authentication proxy deployed is secure to be open to the internet, so if necessary it is acceptable to run
clusters without restricting the range. However, where possible, the range should be restricted to those IP ranges
where your administrators will access the cluster from.
`
}

// CheckPlan checks a plan for compliance with this rule
func (p *AuthProxyIPRangeRule) CheckPlan(ctx context.Context, client client.Client, target *configv1.Plan) (*securityv1.SecurityScanRuleResult, error) {
	result := &securityv1.SecurityScanRuleResult{
		RuleCode: p.Code(),
		Status:   securityv1.Warning,
	}

	var config map[string]interface{}
	if err := json.Unmarshal(target.Spec.Configuration.Raw, &config); err != nil {
		result.Message = "Could not check plan as plan JSON invalid"
		return result, err
	}

	ipRanges, ok := config["authProxyAllowedIPs"].([]interface{})
	if !ok {
		result.Message = "Could not check plan as plan authProxyAllowedIPs not an array"
		return result, fmt.Errorf("authProxyAllowedIPs not an array")
	}

	if len(ipRanges) == 0 {
		result.Message = "No Auth Proxy IP ranges specified by plan"
		return result, nil
	}

	for _, r := range ipRanges {
		rangeStr, ok := r.(string)
		if !ok {
			result.Message = fmt.Sprintf("Range %v not a string, can't check", r)
			return result, fmt.Errorf("Invalid plan value: %s", result.Message)
		}
		bits := strings.Split(rangeStr, "/")
		if len(bits) != 2 {
			result.Message = fmt.Sprintf("Range %s not in format 0.0.0.0/0, can't check", rangeStr)
			return result, fmt.Errorf("Invalid plan value: %s", result.Message)
		}
		net, err := strconv.Atoi(bits[1])
		if err != nil {
			result.Message = fmt.Sprintf("Range %s not in format 0.0.0.0/0, can't check", rangeStr)
			return result, fmt.Errorf("Invalid plan value: %s", result.Message)
		}
		if net < 16 {
			result.Message = fmt.Sprintf("Range %s specifies wider than /16", rangeStr)
			return result, nil
		}
	}

	result.Message = "All ranges specified checked and compliant"
	result.Status = securityv1.Compliant

	return result, nil
}
