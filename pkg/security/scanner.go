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
	"sync"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// New returns an initialized implementation of a Scanner
func New() Scanner {
	scanner := NewEmpty()

	// Register default built-in security rules. Further rules can be added using
	// RegisterRule if required.
	scanner.RegisterRule(&AuthProxyIPRangeRule{})

	return scanner
}

// NewEmpty returns an initialized implementation of scanner without any default
// rules added.
func NewEmpty() Scanner {
	scanner := scannerImpl{
		rulesLock: &sync.RWMutex{},
	}

	// Register default built-in security rules. Further rules can be added using
	// RegisterRule if required.
	scanner.RegisterRule(&AuthProxyIPRangeRule{})

	return &scanner
}

type scannerImpl struct {
	rulesLock *sync.RWMutex
	rules     []Rule
}

func (s *scannerImpl) RegisterRule(rule Rule) {
	s.rulesLock.Lock()
	defer s.rulesLock.Unlock()

	s.rules = append(s.rules, rule)
}

func (s *scannerImpl) GetRules() []Rule {
	s.rulesLock.RLock()
	defer s.rulesLock.RUnlock()

	return s.rules
}

func (s *scannerImpl) GetRule(code string) Rule {
	s.rulesLock.RLock()
	defer s.rulesLock.RUnlock()

	for _, rule := range s.rules {
		if rule.Code() == code {
			return rule
		}
	}
	return nil
}

func (s *scannerImpl) ScanPlan(target *configv1.Plan) *securityv1.ScanResult {
	return s.scanRules(target.TypeMeta, target.ObjectMeta, func(rule Rule) (bool, securityv1.RuleResult) {
		// Apply the rule if it implements PlanRule interface:
		pr, applicable := rule.(PlanRule)
		if !applicable {
			return false, securityv1.RuleResult{}
		}
		return true, pr.CheckPlan(target)
	})
}

func (s *scannerImpl) ScanCluster(target *clustersv1.Cluster) *securityv1.ScanResult {
	return s.scanRules(target.TypeMeta, target.ObjectMeta, func(rule Rule) (bool, securityv1.RuleResult) {
		// Apply the rule if it implements ClusterRule interface:
		cr, applicable := rule.(ClusterRule)
		if !applicable {
			return false, securityv1.RuleResult{}
		}
		return true, cr.CheckCluster(target)
	})
}

func (s *scannerImpl) scanRules(typeMeta metav1.TypeMeta, objMeta metav1.ObjectMeta, scan func(Rule) (bool, securityv1.RuleResult)) *securityv1.ScanResult {
	result := securityv1.ScanResult{
		Spec: securityv1.ScanResultSpec{
			OverallStatus:      securityv1.Compliant,
			ResourceAPIVersion: typeMeta.APIVersion,
			ResourceKind:       typeMeta.Kind,
			ResourceNamespace:  objMeta.Namespace,
			ResourceName:       objMeta.Name,
			CheckedAt:          metav1.NewTime(time.Now()),
		},
	}

	s.rulesLock.RLock()
	defer s.rulesLock.RUnlock()

	for _, rule := range s.rules {
		applicable, ruleResult := scan(rule)
		if applicable {
			ruleResult.CheckedAt = metav1.NewTime(time.Now())
			result.Spec.Results = append(result.Spec.Results, ruleResult)
			setOverallResult(&result, &ruleResult)
		}
	}

	return &result
}

// setOverallResult sets the result on the scan result to that of the rule result only if the
// rule result is WORSE than the overall result, otherwise it leaves it as is.
func setOverallResult(result *securityv1.ScanResult, ruleResult *securityv1.RuleResult) {
	// Downgrade overall result if applicable
	if result.Spec.OverallStatus == securityv1.Compliant && (ruleResult.Status == securityv1.Warning || ruleResult.Status == securityv1.Failure) {
		result.Spec.OverallStatus = ruleResult.Status
	} else if result.Spec.OverallStatus == securityv1.Warning && ruleResult.Status == securityv1.Failure {
		result.Spec.OverallStatus = ruleResult.Status
	}
}
