/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	"fmt"
)

// Team is the contract to a team
type Team interface {
	// Allocations returns the team allocation interface
	Allocations() Allocations
	// Cloud returns the cloud providers
	Cloud() Cloud
	// Clusters returns the teams clusters
	Clusters() Clusters
	// Members returns the team members interface
	Members() TeamMembers
	// Namespace is the name for a team
	Namespace() string
}

// tmImpl is a team interface
type tmImpl struct {
	*hubImpl
	// team is the name of the team
	team string
}

// Allocations return an interface to the team allocations
func (t tmImpl) Allocations() Allocations {
	return &acaImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t tmImpl) Cloud() Cloud {
	return &cloudImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t tmImpl) Clusters() Clusters {
	return &clsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// Members returns the team members interface
func (t tmImpl) Members() TeamMembers {
	return &tmsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// Namespace returns the kubernetes namespace for a team
// @TBD at the moment we are saying all teams are located in a single namespace
// which is prefixed with team-<name>
func (t tmImpl) Namespace() string {
	return fmt.Sprintf(t.team)
}
