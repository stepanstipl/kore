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
	"fmt"
	"strings"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/utils"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

type emrule struct {
	name, description, code, status string
	compiler                        *ast.Compiler
}

// EmbeddedRules uses OPA ruleset to evaluate the assets
type EmbeddedRules struct {
	rules []Rule
}

// NewEmbeddedRules reads the rules
func NewEmbeddedRules() (*EmbeddedRules, error) {
	var list []Rule
	// @step: attempt to load the assets
	for _, x := range AssetNames() {
		policy, err := Asset(x)
		if err != nil {
			log.WithError(err).Error("trying to load the assets from file")

			return nil, err
		}
		policies := make([]map[string]interface{}, 0)

		err = yaml.Unmarshal(policy, &policies)
		if err != nil {
			log.WithField("policy", x).WithError(err).Error("trying to decode the policies")

			return nil, err
		}

		for _, j := range policies {
			rule := j["rule"].(string)
			compiler, err := ast.CompileModules(map[string]string{
				"security.rego": rule,
			})
			if err != nil {
				log.WithError(err).WithField("policy", x).Error("trying to compile the policy")

				return nil, err
			}

			list = append(list, &emrule{
				name:        j["name"].(string),
				code:        j["code"].(string),
				description: j["description"].(string),
				status:      j["status"].(string),
				compiler:    compiler,
			})
		}
	}

	return &EmbeddedRules{rules: list}, nil
}

// List returns a list of rules from the assets
func (r *EmbeddedRules) List() []Rule {
	return r.rules
}

func (e emrule) Code() string {
	return e.code
}

func (e emrule) Name() string {
	return e.name
}

func (e emrule) Description() string {
	return e.description
}

func (e *emrule) performCheck(ctx context.Context, client client.Client, resource runtime.Object) (*securityv1.SecurityScanRuleResult, error) {
	rego := rego.New(
		rego.Query("data.security"),
		rego.Compiler(e.compiler),
		rego.Input(resource),
	)
	rs, err := rego.Eval(ctx)
	if err != nil {
		log.WithError(err).Error("trying to evaluate the rule")

		return nil, err
	}
	if len(rs) == 0 {
		return nil, errors.New("no result found in rule evaluation")
	}

	// @step: process the evaluation
	v := rs[0].Expressions[0].Value
	values, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response from rule %s", e.Code())
	}

	v, found := values["msg"]
	if !found {
		return nil, fmt.Errorf("no default message found in rule %s", e.Code())
	}

	// @step: get the default message
	message, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("default message from rule is not a string")
	}
	status := securityv1.Compliant
	var reason string

	switch {
	case utils.IsEqualType(resource, &configv1.Plan{}):
		v, found := values["plan"].([]interface{})
		if found && len(v) > 0 {
			status = securityv1.RuleStatus(strings.Title(e.status))
			reason = fmt.Sprintf("%s", v[0])
		}
	case utils.IsEqualType(resource, &clustersv1.Cluster{}):
		v, found := values["cluster"].([]interface{})
		if found && len(v) > 0 {
			status = securityv1.RuleStatus(strings.Title(e.status))
			reason = fmt.Sprintf("%s", v[0])
		}
	}
	if reason != "" {
		message = reason
	}

	return &securityv1.SecurityScanRuleResult{
		CheckedAt: metav1.NewTime(time.Now()),
		RuleCode:  e.Code(),
		Message:   message,
		Status:    status,
	}, nil
}

func (e *emrule) CheckPlan(ctx context.Context, client client.Client, target *configv1.Plan) (*securityv1.SecurityScanRuleResult, error) {
	return e.performCheck(ctx, client, target)
}

func (e *emrule) CheckCluster(ctx context.Context, client client.Client, target *clustersv1.Cluster) (*securityv1.SecurityScanRuleResult, error) {
	return e.performCheck(ctx, client, target)
}
