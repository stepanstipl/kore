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

var (
	// List provides a default list
	List ListFuncs
	// Filter providers a filter
	Filter ListFuncs
)

// ListFuncs provides options for listing resources
type ListFuncs struct{}

// WithID sets the id
func (q ListFuncs) WithID(id uint64) ListFunc {
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

// WithProviders sets the provider name
func (q ListFuncs) WithProviders(v []string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["provider.names"] = v
	}
}

// WithProviderToken sets the provider token
func (q ListFuncs) WithProviderToken(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["provider.token"] = v
	}
}

// WithDuration searches for a duration
func (q ListFuncs) WithDuration(since time.Duration) ListFunc {
	return func(o *ListOptions) {
		o.Fields["duration"] = since
	}
}

// WithTeam sets the team
func (q ListFuncs) WithTeam(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["team"] = v
	}
}

// WithStatus sets the status
func (q ListFuncs) WithStatus(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["status"] = v
	}
}

// NotTeam sets the team
func (q ListFuncs) NotTeam(v ...string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["team.not"] = append([]string{}, v...)
	}
}

// WithTeamNotNull ensures the team is there
func (q ListFuncs) WithTeamNotNull() ListFunc {
	return func(o *ListOptions) {
		o.Fields["teams.not_null"] = true
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

// WithName sets the name
func (q ListFuncs) WithName(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["name"] = v
	}
}

// WithResourceName sets the resource name
func (q ListFuncs) WithResourceName(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["resource.name"] = v
	}
}

// WithNamespace sets the namespace
func (q ListFuncs) WithNamespace(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["namespace"] = v
	}
}

// WithIdentity sets the API version, kind, namespace and name
func (q ListFuncs) WithIdentity(group string, version string, kind string, namespace string, name string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["group"] = group
		o.Fields["version"] = version
		o.Fields["kind"] = kind
		o.Fields["namespace"] = namespace
		o.Fields["name"] = name
	}
}

// WithResourceGroup set the resource group
func (q ListFuncs) WithResourceGroup(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["group"] = v
	}
}

// WithResourceVersion set the resource group
func (q ListFuncs) WithResourceVersion(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["version"] = v
	}
}

// WithResourceKind set the resource group
func (q ListFuncs) WithResourceKind(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["kind"] = v
	}
}

// WithAlertLatest sets the latest flag on the search
func (q ListFuncs) WithAlertLatest() ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.latest"] = true
	}
}

// WithAlertHistory sets the history on the search
func (q ListFuncs) WithAlertHistory(v int) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.history"] = v
	}
}

// WithAlertSource sets the source on the search
func (q ListFuncs) WithAlertSource(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.source"] = v
	}
}

// WithAlertFingerprint sets the source on the search
func (q ListFuncs) WithAlertFingerprint(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.fingerprint"] = v
	}
}

// WithAlertUID sets the source on the search
func (q ListFuncs) WithAlertUID(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.uid"] = v
	}
}

// WithAlertStatus sets the status flag on the search
func (q ListFuncs) WithAlertStatus(v []string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.status"] = v
	}
}

// WithAlertLabels sets the status flag on the search
func (q ListFuncs) WithAlertLabels(v []string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.labels"] = v
	}
}

// WithAlertSeverity sets the severity flag on the search
func (q ListFuncs) WithAlertSeverity(v string) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.severity"] = v
	}
}

// WithRuleID sets the severity flag on the search
func (q ListFuncs) WithRuleID(v uint64) ListFunc {
	return func(o *ListOptions) {
		o.Fields["alert.rule_id"] = v
	}
}
