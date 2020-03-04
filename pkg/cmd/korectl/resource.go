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

package korectl

import (
	"github.com/go-openapi/inflect"
)

// resourceConfig stores custom resource CLI configurations
type resourceConfig struct {
	APIEndpoint    string
	RequiredParams []string
	Columns        []string
}

func getResourceConfig(name string, team string) (resourceConfig, error) {
	if team == "" {
		if config, ok := globalResourceConfigs[inflect.Singularize(name)]; ok {
			return config, nil
		}
	}

	if config, ok := teamResourceConfigs[inflect.Singularize(name)]; ok {
		return config, nil
	}

	return resourceConfig{
		APIEndpoint:    "/teams/{team}/" + name,
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	}, nil
}

var globalResourceConfigs = map[string]resourceConfig{
	"team": {
		APIEndpoint: "/teams",
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.description"),
		},
	},
	"user": {
		APIEndpoint: "/users",
		Columns: []string{
			Column("Username", ".metadata.name"),
			Column("Email", ".spec.email"),
			Column("Disabled", ".spec.disabled"),
		},
	},
	"plan": {
		APIEndpoint: "/plans",
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Description", ".spec.description"),
			Column("Summary", ".spec.summary"),
		},
	},
	"audit-event": {
		APIEndpoint: "/audit",
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	},
}

var teamResourceConfigs = map[string]resourceConfig{
	"team-member": {
		APIEndpoint:    "/teams/{team}/members",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Username", "."),
		},
	},
	"allocation": {
		APIEndpoint:    "/teams/{team}/allocations",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.summary"),
			Column("Resource", ".spec.resource.kind"),
		},
	},
	"cluster": {
		APIEndpoint:    "/teams/{team}/clusters",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Provider", ".spec.provider.group"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	},
	"namespaceclaim": {
		APIEndpoint:    "/teams/{team}/namespaceclaims",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Namespace", ".spec.name"),
			Column("Cluster", ".spec.cluster.name"),
			Column("Status", ".status.status"),
		},
	},
	"audit-event": {
		APIEndpoint:    "/teams/{team}/audit",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	},
}
