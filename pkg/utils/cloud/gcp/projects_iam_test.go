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
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func makeTestPolicy() *cloudresourcemanager.Policy {
	return &cloudresourcemanager.Policy{
		Bindings: []*cloudresourcemanager.Binding{
			{
				Members: []string{"serviceAccount:service-633278442909@compute-system.iam.gserviceaccount.com"},
				Role:    "roles/compute.serviceAgent",
			},
			{
				Members: []string{"serviceAccount:service-633278442909@container-engine-robot.iam.gserviceaccount.com"},
				Role:    "roles/container.serviceAgent",
			},
			{
				Members: []string{
					"serviceAccount:633278442909-compute@developer.gserviceaccount.com",
					"serviceAccount:633278442909@cloudservices.gserviceaccount.com",
					"serviceAccount:service-633278442909@containerregistry.iam.gserviceaccount.com",
				},
				Role: "roles/editor",
			},
			{
				Members: []string{
					"serviceAccount:kore-admin@kore-admin.iam.gserviceaccount.com",
					"user: test@appvia.io",
				},
				Role: "roles/owner",
			},
		},
	}
}

func TestAddBindingsToProjectPolicy(t *testing.T) {
	cases := []struct {
		Bindings []*cloudresourcemanager.Binding
		Expected func() []*cloudresourcemanager.Binding
	}{
		{
			Bindings: []*cloudresourcemanager.Binding{
				{
					Role:    "roles/owner",
					Members: []string{"serviceAccount:test@kore-admin.iam.gserviceaccount.com"},
				},
			},
			Expected: func() []*cloudresourcemanager.Binding {
				p := makeTestPolicy()
				p.Bindings[3].Members = append(p.Bindings[3].Members, "serviceAccount:test@kore-admin.iam.gserviceaccount.com")

				return p.Bindings
			},
		},
		{
			Bindings: []*cloudresourcemanager.Binding{
				{
					Role:    "roles/new",
					Members: []string{"serviceAccount:test"},
				},
			},
			Expected: func() []*cloudresourcemanager.Binding {
				p := makeTestPolicy()
				p.Bindings = append(p.Bindings, &cloudresourcemanager.Binding{
					Role:    "roles/new",
					Members: []string{"serviceAccount:test"},
				})

				return p.Bindings
			},
		},
	}
	for i, c := range cases {
		policy := makeTestPolicy()
		require.NoError(t, AddBindingsToProjectPolicy(policy, c.Bindings))
		if !assert.Equal(t, c.Expected(), policy.Bindings, "case %d not as expected", i) {
			fmt.Println("Expected:", spew.Sdump(c.Expected()))
			fmt.Println("Got:", spew.Sdump(policy.Bindings))
		}
	}
}
