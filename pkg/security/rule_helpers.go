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
