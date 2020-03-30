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

import (
	"strings"

	"github.com/appvia/kore/pkg/utils"
)

type resourceConfig struct {
	Name     string
	IsGlobal bool
	IsTeam   bool
	Columns  []string
}

func getResourceConfig(name string) resourceConfig {
	name = strings.ToLower(name)
	switch name {
	case "eks", "ekss":
		name = "eks"
	default:
		name = utils.Singularize(name)
	}

	if config, ok := resourceConfigs[name]; ok {
		return config
	}

	// TODO: we don't have a way to validate whether a resource exists yet, so we generate a team resource configuration dynamically
	return resourceConfig{
		Name:   utils.Pluralize(name),
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
	"audit": {
		Name:     "audit",
		IsGlobal: true,
		IsTeam:   true,
		Columns: []string{
			Column("Time", "spec.createdAt"),
			Column("Operation", "spec.operation"),
			Column("URI", "spec.resourceURI"),
			Column("User", "spec.user"),
			Column("Team", "spec.team"),
			Column("Result", "spec.responseCode"),
		},
	},
	"kubernetes": {
		Name:   "kubernetes",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Provider", "spec.provider.group"),
			Column("Endpoint", "status.endpoint"),
			Column("Status", "status.status"),
		},
	},
	"cluster": {
		Name:   "clusters",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Kind", "spec.kind"),
			Column("API Endpoint", "status.apiEndpoint"),
			Column("Auth Proxy Endpoint", "status.authProxyEndpoint"),
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
	"gkecredential": {
		Name:   "gkecredentials",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Project", "spec.project"),
			Column("Status", "status.status"),
			Column("Verified", "status.verified"),
		},
	},
	"eks": {
		Name:   "ekss",
		IsTeam: true,
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Region", "spec.region"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	},
	"ekscredential": {
		Name:   "ekscredentials",
		IsTeam: true,
		Columns: []string{
			Column("Name", "metadata.name"),
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
			Column("Summary", "spec.summary"),
			Column("Description", "spec.description"),
		},
	},
	"planpolicy": {
		Name:     "planpolicies",
		IsGlobal: true,
		Columns: []string{
			Column("Resource", "metadata.name"),
			Column("Summary", "spec.summary"),
			Column("Description", "spec.description"),
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
			Column("Status", "status.status"),
		},
	},
	"team": {
		Name:     "teams",
		IsGlobal: true,
		Columns: []string{
			Column("Name", "metadata.name"),
			Column("Description", "spec.description"),
			Column("Status", "status.status"),
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
