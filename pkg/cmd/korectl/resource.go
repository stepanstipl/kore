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

package korectl

import "github.com/go-openapi/inflect"

type resourceConfig struct {
	Name     string
	IsGlobal bool
	IsTeam   bool
	Columns  []string
}

func getResourceConfig(name string) resourceConfig {
	if config, ok := resourceConfigs[inflect.Singularize(name)]; ok {
		return config
	}

	// TODO: we don't have a way to validate whether a resource exists yet, so we generate a team resource configuration dynamically
	return resourceConfig{
		Name:   name,
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
		},
	}
}

var resourceConfigs = map[string]resourceConfig{
	"allocation": {
		Name:   "allocations",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Description", "spec.summary"),
			Column("Owned By", "metadata.namespace"),
			Column("Resource", "spec.resource.kind"),
			Column("Status", "status.status"),
		},
	},
	"audit-event": {
		Name:     "audit",
		IsGlobal: true,
		IsTeam:   true,
		Columns: []string{
			Column("Time", "spec.createdAt"),
			Column("Resource", "spec.resource"),
			Column("URI", "spec.resourceURI"),
			Column("Operation", "spec.operation"),
			Column("User", "spec.user"),
			Column("Team", "spec.team"),
			Column("Result", "spec.responseCode"),
		},
	},
	"cluster": {
		Name:   "clusters",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Provider", "spec.provider.group"),
			Column("Endpoint", "status.endpoint"),
			Column("Status", "status.status"),
		},
	},
	"gke": {
		Name:   "gkes",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Region", "spec.region"),
			Column("Endpoint", "status.endpoint"),
			Column("Status", "status.status"),
		},
	},
	"member": {
		Name:   "members",
		IsTeam: true,
		Columns: []string{
			Column("Username", ""),
		},
	},
	"namespaceclaim": {
		Name:   "namespaceclaims",
		IsTeam: true,
		Columns: []string{
			Column("Resource", "metadata.name"),
			Column("Namespace", "spec.name"),
			Column("Cluster", "spec.cluster.name"),
			Column("Status", "status.status"),
		},
	},
	"plan": {
		Name:     "plans",
		IsGlobal: true,
		Columns: []string{
			Column("Resource", "metadata.name"),
			Column("Description", "spec.description"),
			Column("Summary", "spec.summary"),
		},
	},
	"organization": {
		Name:   "organizations",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Status", "status.status"),
		},
	},
	"projectclaim": {
		Name:   "projectclaims",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Organization", "spec.organization.name."),
			Column("Owned By", "spec.organization.namespace"),
			Column("Status", "status.status"),
		},
	},
	"secret": {
		Name:   "secrets",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Type", "spec.type"),
			Column("Description", "spec.description"),
			Column("Verified", "status.verified"),
		},
	},
	"team": {
		Name:     "teams",
		IsGlobal: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Description", "spec.description"),
		},
	},
	"user": {
		Name:     "users",
		IsGlobal: true,
		Columns: []string{
			Column("Username", "metadata.name"),
			Column("Email", "spec.email"),
			Column("Disabled", "spec.disabled"),
		},
	},
}
