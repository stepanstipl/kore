/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package hub

// Cloud returns a collection of cloud providers
type Cloud interface {
	// GKE returms the GKE interface
	GKE() GKE
	// GKE
}

type cloudImpl struct {
	*hubImpl
	// team is the requesting team
	team string
}

// GKE returns a gke interface
func (c *cloudImpl) GKE() GKE {
	return &gkeImpl{cloudImpl: c, team: c.team}
}
