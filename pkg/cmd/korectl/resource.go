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
	} else {
		// TODO: we don't have a way to validate whether a resource exists yet, so we generate a team resource configuration dynamically
		return resourceConfig{
			Name:   name,
			IsTeam: true,
			Columns: []string{
				Column("Name", ".metadata.name"),
			},
		}
	}
}

var resourceConfigs = map[string]resourceConfig{
	"allocation": {
		Name:   "allocations",
		IsTeam: true,
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.summary"),
			Column("Resource", ".spec.resource.kind"),
		},
	},
	"audit-event": {
		Name:     "audit",
		IsGlobal: true,
		IsTeam:   true,
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	},
	"cluster": {
		Name:   "clusters",
		IsTeam: true,
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Provider", ".spec.provider.group"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	},
	"gke": {
		Name:   "gkes",
		IsTeam: true,
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	},
	"member": {
		Name:   "members",
		IsTeam: true,
		Columns: []string{
			Column("Username", "."),
		},
	},
	"namespaceclaim": {
		Name:   "namespaceclaims",
		IsTeam: true,
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Namespace", ".spec.name"),
			Column("Cluster", ".spec.cluster.name"),
			Column("Status", ".status.status"),
		},
	},
	"plan": {
		Name:     "plans",
		IsGlobal: true,
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Description", ".spec.description"),
			Column("Summary", ".spec.summary"),
		},
	},
	"team": {
		Name:     "teams",
		IsGlobal: true,
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.description"),
		},
	},
	"user": {
		Name:     "users",
		IsGlobal: true,
		Columns: []string{
			Column("Username", ".metadata.name"),
			Column("Email", ".spec.email"),
			Column("Disabled", ".spec.disabled"),
		},
	},
}
