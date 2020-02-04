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

package users

// ListFunc are terms to search on
type ListFunc func(*ListOptions)

// ListOptions defines the where clause of the search
type ListOptions struct {
	Fields map[string]interface{}
}

// NewListOptions returns a list options
func NewListOptions() *ListOptions {
	return &ListOptions{Fields: make(map[string]interface{})}
}

// ApplyListOptions is responsible for applying the terms
func ApplyListOptions(v ...ListFunc) *ListOptions {
	o := NewListOptions()

	for _, x := range v {
		x(o)
	}

	return o
}

// Has checks if a field is set
func (l *ListOptions) Has(k string) bool {
	_, found := l.Fields[k]

	return found
}

// GetString returns a string from the fields
func (l *ListOptions) GetString(k string) string {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// GetStringSlice returns a string slice
func (l *ListOptions) GetStringSlice(k string) []string {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.([]string); ok {
			return s
		}

	}
	return []string{}
}

// GetInt returns an int from the fields
func (l *ListOptions) GetInt(k string) int {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.(int); ok {
			return s
		}

		return 0
	}

	return 0
}

// GetBool returns the boolean
func (l *ListOptions) GetBool(k string) bool {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.(bool); ok {
			return s
		}

		return false
	}

	return false
}

// HasID checks the id
func (l *ListOptions) HasID() bool {
	return l.Has("id")
}

// HasNotNames checks for not names
func (l *ListOptions) HasNotNames() bool {
	return l.Has("not.names")
}

// HasName checks the name
func (l *ListOptions) HasName() bool {
	return l.Has("name")
}

// HasProvider checks the name
func (l *ListOptions) HasProvider() bool {
	return l.Has("provider.name")
}

// HasProviderToken checks the name
func (l *ListOptions) HasProviderToken() bool {
	return l.Has("provider.token")
}

// HasTeam checks the team
func (l *ListOptions) HasTeam() bool {
	return l.Has("team")
}

// HasNotTeam checks the team
func (l *ListOptions) HasNotTeam() bool {
	return l.Has("team.not")
}

// HasTeams checks the team
func (l *ListOptions) HasTeams() bool {
	return l.Has("teams")
}

// HasTeamID checks for a team id
func (l *ListOptions) HasTeamID() bool {
	return l.Has("team.id")
}

// HasEnabled checks the enabled
func (l *ListOptions) HasEnabled() bool {
	return l.Has("enabled")
}

// HasDisabled checks the disable
func (l *ListOptions) HasDisabled() bool {
	return l.Has("disabled")
}

// HasUser checks the user
func (l *ListOptions) HasUser() bool {
	return l.Has("user")
}

// GetID gets the id
func (l *ListOptions) GetID() int {
	return l.GetInt("id")
}

// GetNotNames gets the name
func (l *ListOptions) GetNotNames() []string {
	return l.GetStringSlice("not.names")
}

// GetName gets the name
func (l *ListOptions) GetName() string {
	return l.GetString("name")
}

// GetProvider gets the name
func (l *ListOptions) GetProvider() string {
	return l.GetString("provider.name")
}

// GetProviderToken gets the provider token
func (l *ListOptions) GetProviderToken() string {
	return l.GetString("provider.token")
}

// GetEnabled gets the enabled
func (l *ListOptions) GetEnabled() bool {
	return l.GetBool("enabled")
}

// GetDisabled gets the enabled
func (l *ListOptions) GetDisabled() bool {
	return l.GetBool("disabled")
}

// GetTeam gets the team
func (l *ListOptions) GetTeam() string {
	return l.GetString("team")
}

// GetNotTeam gets the team
func (l *ListOptions) GetNotTeam() []string {
	return l.GetStringSlice("team.not")
}

// GetTeams gets the team
func (l *ListOptions) GetTeams() []string {
	return l.GetStringSlice("teams")
}

// GetTeamID gets the team id
func (l *ListOptions) GetTeamID() int {
	return l.GetInt("team.id")
}

// GetUser gets the user
func (l *ListOptions) GetUser() string {
	return l.GetString("user")
}
