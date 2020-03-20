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

package kore

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
	// NamespaceClaims returns the the interface
	NamespaceClaims() NamespaceClaims
	// Secrets returns the secret interface
	Secrets() Secrets
}

// tmImpl is a team interface
type tmImpl struct {
	*hubImpl
	// team is the name of the team
	team string
}

// Allocations return an interface to the team allocations
func (t *tmImpl) Allocations() Allocations {
	return &acaImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t *tmImpl) Cloud() Cloud {
	return &cloudImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

func (t *tmImpl) Clusters() Clusters {
	return &clsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// Members returns the team members interface
func (t *tmImpl) Members() TeamMembers {
	return &tmsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// NamespaceClaims returns a namespace claim interface
func (t *tmImpl) NamespaceClaims() NamespaceClaims {
	return &nsImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}

// Secrets returns a secrets interface
func (t *tmImpl) Secrets() Secrets {
	return &secretImpl{
		hubImpl: t.hubImpl,
		team:    t.team,
	}
}
