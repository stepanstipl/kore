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

package users

var (
	// List provides a default list
	List ListFuncs
	// Filter providers a filter
	Filter ListFuncs
)

// ListFuncs provides options for listing resources
type ListFuncs struct{}

// WithID sets the id
func (q ListFuncs) WithID(id int) ListFunc {
	return func(o *ListOptions) {
		o.Fields["id"] = id
	}
}

// NotNames set the inverse of names
func (q ListFuncs) NotNames(v ...string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["not.names"] = v
	}
}

// WithProvider sets the provider name
func (q ListFuncs) WithProvider(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["provider.name"] = v
	}
}

// WithProviderToken sets the provider token
func (q ListFunc) WithProviderToken(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["provider.token"] = v
	}
}

// WithTeam sets the team
func (q ListFuncs) WithTeam(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["team"] = v
	}
}

// NotTeam sets the team
func (q ListFuncs) NotTeam(v ...string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["team.not"] = append([]string{}, v...)
	}
}

// WithTeams sets the team
func (q ListFuncs) WithTeams(v ...string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["teams"] = append([]string{}, v...)
	}
}

// WithTeamID sets the team id
func (q ListFuncs) WithTeamID(v uint64) ListFunc {
	return func(o *ListOptions) {
		o.Fields["team.id"] = int(v)
	}
}

// WithDisabled sets the disabled
func (q ListFuncs) WithDisabled(v bool) ListFunc {
	return func(o *ListOptions) {
		o.Fields["disabled"] = v
	}
}

// WithEnabled sets the enabled
func (q ListFuncs) WithEnabled(v bool) ListFunc {
	return func(o *ListOptions) {
		o.Fields["enabled"] = v
	}
}

// WithUser sets the user
func (q ListFuncs) WithUser(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["user"] = v
	}
}
