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

package audit

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

// WithUser sets the user
func (q ListFuncs) WithUser(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["user"] = v
	}
}
