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
	Name            string
	APIResourceName string
	Columns         []string
}

type resourceConfigMap map[string]resourceConfig

func (r resourceConfigMap) Get(name string) resourceConfig {
	if config, ok := r[name]; ok {
		return config
	}

	return resourceConfig{
		Name:            name,
		APIResourceName: name,
		Columns: []string{
			Column("Name", ".metadata.name"),
		},
	}
}

// @question: we might wanna consider renaming these to printers
var (
	teamResourceConfig = resourceConfig{
		Name:            "team",
		APIResourceName: "teams",
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.description"),
		},
	}

	allocationResourceConfig = resourceConfig{
		Name:            "allocation",
		APIResourceName: "allocations",
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Description", ".spec.summary"),
			Column("Resource", ".spec.resource.kind"),
		},
	}
	clusterResourceConfig = resourceConfig{
		Name:            "cluster",
		APIResourceName: "clusters",
		Columns: []string{
			Column("Name", ".metadata.name"),
			Column("Provider", ".spec.provider.group"),
			Column("Endpoint", ".status.endpoint"),
			Column("Status", ".status.status"),
		},
	}
	namespaceResourceConfig = resourceConfig{
		Name:            "namespaceclaim",
		APIResourceName: "namespaceclaims",
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Namespace", ".spec.name"),
			Column("Cluster", ".spec.cluster.name"),
			Column("Status", ".status.status"),
		},
	}
	planResourceConfig = resourceConfig{
		Name:            "plan",
		APIResourceName: "plans",
		Columns: []string{
			Column("Resource", ".metadata.name"),
			Column("Description", ".spec.description"),
			Column("Summary", ".spec.summary"),
		},
	}
)

var resourceConfigs = resourceConfigMap{
	"allocation":      allocationResourceConfig,
	"allocations":     allocationResourceConfig,
	"cluster":         clusterResourceConfig,
	"clusters":        clusterResourceConfig,
	"namespaceclaim":  namespaceResourceConfig,
	"namespaceclaims": namespaceResourceConfig,
	"plan":            planResourceConfig,
	"plans":           planResourceConfig,
	"team":            teamResourceConfig,
	"teams":           teamResourceConfig,
}
