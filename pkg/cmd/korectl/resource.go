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

// resourceConfig stores custom resource CLI configurations
type resourceConfig struct {
	APIEndpoint    string
	RequiredParams []string
	Columns        []string
}

type resourceConfigMap map[string]*resourceConfig

func (r resourceConfigMap) Get(name string) *resourceConfig {
	if config, ok := r[name]; ok {
		return config
	}

	return &resourceConfig{
		APIEndpoint:    "/teams/{team}/" + name,
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	}
}

var resourceConfigs = resourceConfigMap{
	"teams": &resourceConfig{
		APIEndpoint: "/teams",
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.description"),
		},
	},

	"users": &resourceConfig{
		APIEndpoint: "/users",
		Columns: []string{
			Column("Username", ".metadata.name"),
			Column("Email", ".spec.email"),
			Column("Disabled", ".spec.disabled"),
		},
	},
	"plan": &resourceConfig{
		APIEndpoint: "/plans",
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Description", ".spec.description"),
			Column("Summary", ".spec.summary"),
		},
	},
	"team-members": &resourceConfig{
		APIEndpoint:    "/teams/{team}/members",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Username", "."),
		},
	},

	"allocation": &resourceConfig{
		APIEndpoint:    "/teams/{team}/allocations",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.summary"),
			Column("Resource", ".spec.resource.kind"),
		},
	},

	"cluster": &resourceConfig{
		APIEndpoint:    "/teams/{team}/clusters",
		RequiredParams: []string{"team"},
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Provider", ".spec.provider.group"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	},
}
