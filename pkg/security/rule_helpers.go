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
	"fmt"

	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/tidwall/gjson"
)

// RuleApplies returns a list of strings indicating the resource kinds that a rule
// can be applied to.
func RuleApplies(rule Rule) []string {
	var appliesTo []string

	if _, applies := rule.(PlanRule); applies {
		appliesTo = append(appliesTo, "Plan")
	}
	if _, applies := rule.(ClusterRule); applies {
		appliesTo = append(appliesTo, "Cluster")
	}

	return appliesTo
}

// ValueAsExpected checks the value is as expected
func ValueAsExpected(code, config, field string, expected interface{}, failStatus securityv1.RuleStatus, success, failure string) (*securityv1.SecurityScanRuleResult, error) {
	result := &securityv1.SecurityScanRuleResult{
		RuleCode: code,
		Status:   securityv1.Warning,
	}

	value := gjson.Get(config, field)
	if !value.Exists() {
		result.Message = "Could not check cluster due to invalid JSON"

		return nil, fmt.Errorf("%s parameter does not exist", field)
	}

	if value.Value() != expected {
		result.Status = failStatus
		result.Message = failure

		return result, nil
	}
	result.Status = securityv1.Compliant

	return result, nil
}
