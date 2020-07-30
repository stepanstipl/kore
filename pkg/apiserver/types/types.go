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

package types

// IssuedToken is a minted token from kore
type IssuedToken struct {
	// Token is the actual token
	Token []byte
	// Expires is the time token will expire
	Expires int64
}

// WhoAmI provides a description to who you are
type WhoAmI struct {
	// AuthMethod is the authentication method being used
	AuthMethod string `json:"authMethod,omitempty"`
	// Email is the user email
	Email string `json:"email,omitempty"`
	// Username is your username
	Username string `json:"username,omitempty"`
	// Teams is a collection of teams your in
	Teams []string `json:"teams,omitempty"`
}

// TeamInvitationResponse returns the team
type TeamInvitationResponse struct {
	// Team is the name of team which the user just has been been added to
	Team string `json:"team"`
}

// Health provides an indication of the health of the API.
type Health struct {
	// Healthy is true if the service is healthy.
	Healthy bool `json:"healthy"`
}
