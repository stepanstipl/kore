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

// HasUser checks the username
func (l *ListOptions) HasUser() bool {
	return l.Has("user.name")
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

// HasType checks the type
func (l *ListOptions) HasType() bool {
	return l.Has("audit.type")
}

// GetID gets the id
func (l *ListOptions) GetID() int {
	return l.GetInt("id")
}

// GetTeam gets the team
func (l *ListOptions) GetTeam() string {
	return l.GetString("team")
}

// GetType checks the type
func (l *ListOptions) GetType() string {
	return l.GetString("audit.type")
}

// GetNotTeam gets the team
func (l *ListOptions) GetNotTeam() []string {
	return l.GetStringSlice("team.not")
}

// GetTeams gets the team
func (l *ListOptions) GetTeams() []string {
	return l.GetStringSlice("teams")
}

// GetUser gets the user
func (l *ListOptions) GetUser() string {
	return l.GetString("user.name")
}
