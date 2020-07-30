/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package persistence

import "time"

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

// GetDurationValue returns a duration
func (l *ListOptions) GetDurationValue(k string) time.Duration {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.(time.Duration); ok {
			return s
		}
	}

	return time.Duration(0)
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

// GetUint64 converts the type
func (l *ListOptions) GetUint64(k string) uint64 {
	v, found := l.Fields[k]
	if found {
		if s, ok := v.(uint64); ok {
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

// HasResourceName checks the name
func (l *ListOptions) HasResourceName() bool {
	return l.Has("resource.name")
}

// HasNamespace checks the namespace
func (l *ListOptions) HasNamespace() bool {
	return l.Has("namespace")
}

// HasProvider checks the name
func (l *ListOptions) HasProvider() bool {
	return l.Has("provider.name")
}

// HasProviders checks the name
func (l *ListOptions) HasProviders() bool {
	return l.Has("provider.names")
}

// HasProviderToken checks the name
func (l *ListOptions) HasProviderToken() bool {
	return l.Has("provider.token")
}

// HasDuration checks for a duration
func (l *ListOptions) HasDuration() bool {
	return l.Has("duration")
}

// HasTeam checks the team
func (l *ListOptions) HasTeam() bool {
	return l.Has("team")
}

// HasAlertLatest checks of the alert latest
func (l *ListOptions) HasAlertLatest() bool {
	return l.Has("alert.latest")
}

// HasAlertHistory checks of the alert history
func (l *ListOptions) HasAlertHistory() bool {
	return l.Has("alert.history")
}

// HasAlertSource checks of the alert source
func (l *ListOptions) HasAlertSource() bool {
	return l.Has("alert.source")
}

// HasAlertFingerprint checks of the alert source
func (l *ListOptions) HasAlertFingerprint() bool {
	return l.Has("alert.fingerprint")
}

// HasAlertUID checks of the alert source
func (l *ListOptions) HasAlertUID() bool {
	return l.Has("alert.uid")
}

// HasAlertStatus checks of the alert status
func (l *ListOptions) HasAlertStatus() bool {
	return l.Has("alert.status")
}

// HasAlertLabels checks of the alert status
func (l *ListOptions) HasAlertLabels() bool {
	return l.Has("alert.labels")
}

// HasRuleID checks for a rule id
func (l *ListOptions) HasRuleID() bool {
	return l.Has("alert.rule_id")
}

// HasAlertSeverity checks of the alert severity
func (l *ListOptions) HasAlertSeverity() bool {
	return l.Has("alert.severity")
}

// HasNotTeam checks the team
func (l *ListOptions) HasNotTeam() bool {
	return l.Has("team.not")
}

// HasTeams checks the team
func (l *ListOptions) HasTeams() bool {
	return l.Has("teams")
}

// HasTeamsNotNull checks the team
func (l *ListOptions) HasTeamsNotNull() bool {
	return l.Has("teams.not_null")
}

// HasTeamID checks for a team id
func (l *ListOptions) HasTeamID() bool {
	return l.Has("team.id")
}

// HasVerb checks the type
func (l *ListOptions) HasVerb() bool {
	return l.Has("audit.verb")
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

// HasIdentity checks if all of the identity fields are present
func (l *ListOptions) HasIdentity() bool {
	return l.HasGroup() && l.HasVersion() && l.HasKind() && l.HasNamespace() && l.HasName()
}

// HasGroup checks if group set
func (l *ListOptions) HasGroup() bool {
	return l.Has("group")
}

// HasVersion checks if version set
func (l *ListOptions) HasVersion() bool {
	return l.Has("version")
}

// HasStatus checks if version set
func (l *ListOptions) HasStatus() bool {
	return l.Has("status")
}

// HasKind checks if kind set
func (l *ListOptions) HasKind() bool {
	return l.Has("kind")
}

// GetID gets the id
func (l *ListOptions) GetID() uint64 {
	return l.GetUint64("id")
}

// GetNotNames gets the name
func (l *ListOptions) GetNotNames() []string {
	return l.GetStringSlice("not.names")
}

// GetName gets the name
func (l *ListOptions) GetName() string {
	return l.GetString("name")
}

// GetResourceName gets the name
func (l *ListOptions) GetResourceName() string {
	return l.GetString("resource.name")
}

// GetNamespace gets the namespace
func (l *ListOptions) GetNamespace() string {
	return l.GetString("namespace")
}

// GetGroup gets the group
func (l *ListOptions) GetGroup() string {
	return l.GetString("group")
}

// GetVersion gets the version
func (l *ListOptions) GetVersion() string {
	return l.GetString("version")
}

// GetKind gets the kind
func (l *ListOptions) GetKind() string {
	return l.GetString("kind")
}

// GetProvider gets the name
func (l *ListOptions) GetProvider() string {
	return l.GetString("provider.name")
}

// GetProviders gets the name
func (l *ListOptions) GetProviders() []string {
	return l.GetStringSlice("provider.names")
}

// GetProviderToken gets the provider token
func (l *ListOptions) GetProviderToken() string {
	return l.GetString("provider.token")
}

// GetDuration gets the duration
func (l *ListOptions) GetDuration() time.Duration {
	return l.GetDurationValue("duration")
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

// GetVerb checks the type
func (l *ListOptions) GetVerb() string {
	return l.GetString("audit.verb")
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

// GetStatus gets the user
func (l *ListOptions) GetStatus() string {
	return l.GetString("status")
}

// GetAlertHistory returns the alert history
func (l *ListOptions) GetAlertHistory() int {
	return l.GetInt("alert.history")
}

// GetAlertSource returns the alert source
func (l *ListOptions) GetAlertSource() string {
	return l.GetString("alert.source")
}

// GetAlertFingerprint returns the alert source
func (l *ListOptions) GetAlertFingerprint() string {
	return l.GetString("alert.fingerprint")
}

// GetAlertUID returns the alert source
func (l *ListOptions) GetAlertUID() string {
	return l.GetString("alert.uid")
}

// GetAlertStatus returns the alert status
func (l *ListOptions) GetAlertStatus() []string {
	return l.GetStringSlice("alert.status")
}

// GetAlertLabels returns the alert status
func (l *ListOptions) GetAlertLabels() []string {
	return l.GetStringSlice("alert.labels")
}

// GetAlertSeverity returns the alert severity
func (l *ListOptions) GetAlertSeverity() string {
	return l.GetString("alert.severity")
}

// GetRuleID returns the rule id
func (l *ListOptions) GetRuleID() uint64 {
	return l.GetUint64("alert.rule_id")
}
