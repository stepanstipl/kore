/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
