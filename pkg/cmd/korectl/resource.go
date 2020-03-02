/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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

var teamResourceConfig = resourceConfig{
	Name:            "team",
	APIResourceName: "teams",
	Columns: []string{
		Column("Name", ".metadata.name"),
		Column("Description", ".spec.description"),
	},
}

var resourceConfigs = resourceConfigMap{
	"team":  teamResourceConfig,
	"teams": teamResourceConfig,
}
