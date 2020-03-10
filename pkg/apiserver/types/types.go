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

package types

// WhoAmI provides a description to who you are
type WhoAmI struct {
	// Email is the user email
	Email string `json:"email,omitempty"`
	// Username is your username
	Username string `json:"username,omitempty"`
	// Teams is a collection of teams your in
	Teams []string `json:"teams,omitempty"`
}

type TeamInvitationResponse struct {
	// Team is the name of team which the user just has been been added to
	Team string `json:"team"`
}
