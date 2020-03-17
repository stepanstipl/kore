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

package gcp

import (
	"context"

	"github.com/appvia/kore/pkg/utils"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// AddBindingsToProjectIAM is responsible for adding the bindings to the projects policy
func AddBindingsToProjectIAM(
	ctx context.Context,
	client *cloudresourcemanager.Service,
	bindings []*cloudresourcemanager.Binding,
	project string) error {

	// @step: first we have to retrieve the policy
	policy, err := client.Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		return err
	}

	// @step: merge in the bindings to the policy
	if err := AddBindingsToProjectPolicy(policy, bindings); err != nil {
		return err
	}

	// @step: update the policy
	if _, err := client.Projects.SetIamPolicy(project, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	}).Context(ctx).Do(); err != nil {
		return err
	}

	return nil
}

// RemoveBindingsFromProjectPolicy is used to remove the bindings from the poliy
func RemoveBindingsFromProjectPolicy(policy *cloudresourcemanager.Policy, bindings []*cloudresourcemanager.Binding) error {
	// @step: iterate the bindings and remove the roles and members
	for _, b := range bindings {
		for j := 0; j < len(policy.Bindings); j++ {
			if b.Role == policy.Bindings[j].Role {
				// we need to remove the users
				for i := 0; i < len(policy.Bindings[j].Members); i++ {
					if utils.Contains(policy.Bindings[j].Members[i], b.Members) {
						policy.Bindings[j].Members = append(policy.Bindings[j].Members[:i], policy.Bindings[j].Members[i+1:]...)
					}
				}
				// @check if the if the binding has no members
				if len(policy.Bindings[j].Members) <= 0 {
					policy.Bindings = append(policy.Bindings[:j], policy.Bindings[j+1:]...)
				}
			}
		}
	}

	return nil
}

// AddBindingsToProjectPolicy is used to adding the bindings to the current policy
func AddBindingsToProjectPolicy(policy *cloudresourcemanager.Policy, bindings []*cloudresourcemanager.Binding) error {
	// @step: merge the bindings into the policy
	for _, x := range bindings {
		// @step: we need to check if the role actually exists
		binding, found := func() (*cloudresourcemanager.Binding, bool) {
			for _, b := range policy.Bindings {
				if b.Role == x.Role {
					return b, true
				}
			}

			return nil, false
		}()
		if !found {
			policy.Bindings = append(policy.Bindings, x)

			continue
		}

		// @step: else the role exists and we need to check and merge users
		for _, m := range x.Members {
			if !utils.Contains(m, binding.Members) {
				binding.Members = append(binding.Members, m)
			}
		}
	}

	return nil
}
