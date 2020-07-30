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

package authentication

// ContextKey is the context key
type ContextKey struct{}

// Identity provides the user
type Identity interface {
	// AuthMethod is the method the user logged in with
	AuthMethod() string
	// IsGlobalAdmin checks if the user is a global admin
	IsGlobalAdmin() bool
	// IsMember checks if the user is a member of a team
	IsMember(string) bool
	// Email returns the user email
	Email() string
	// Disabled checks if the user is disabled
	Disabled() bool
	// Username is a unique username for this identity
	Username() string
	// Teams is a list of teams for the user
	Teams() []string
}
