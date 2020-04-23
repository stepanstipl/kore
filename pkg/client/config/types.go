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

package config

const (
	// DefaultKoreConfigPath is the default path for the kore configuration file
	DefaultKoreConfigPath = "${HOME}/.kore/config"
	// DefaultKoreConfigPathEnv is the default name of the env variable for config
	DefaultKoreConfigPathEnv = "KORE_CONFIG"
)

// Config is the configuration for the api
type Config struct {
	// AuthInfos is a collection of credentials
	AuthInfos map[string]*AuthInfo `json:"users,omitempty" yaml:"users"`
	// CurrentProfile is the profile in use at the moment
	CurrentProfile string `json:"current-profile,omitempty" yaml:"current-profile"`
	// Profiles is a collection of profiles
	Profiles map[string]*Profile `json:"profiles,omitempty" yaml:"profiles"`
	// Servers is a collection of api endpoints
	Servers map[string]*Server `json:"servers,omitempty" yaml:"servers"`
	// Version is the version of the configuration
	Version string `json:"version,omitempty" yaml:"version"`
	// FeatureGates shows all feature gates
	FeatureGates map[string]bool `json:"feature-gates,omitempty" yaml:"feature-gates"`
}

// AuthInfo defines a credential to the api endpoint
type AuthInfo struct {
	// BasicAuth defines a basic auth user/pass credential
	BasicAuth *BasicAuth `json:"basic-auth,omitempty" yaml:"basic-auth"`
	// Token is a static token to use
	Token *string `json:"token,omitempty" yaml:"token"`
	// OIDC is credentials from an oauth2 provider
	OIDC *OIDC `json:"oidc,omitempty" yaml:"oidc"`
}

// BasicAuth defines a basic user credential
type BasicAuth struct {
	// Username is the username for authentication
	Username string `json:"username,omitempty" yaml:"username"`
	// Password is the user password
	Password string `json:"password,omitempty" yaml:"password"`
}

// OIDC is the identity within the kore
type OIDC struct {
	// AccessToken is the access token retrieved from kore
	AccessToken string `json:"access-token,omitempty" yaml:"access-token"`
	// AuthorizeURL is the endpoint for the authorize
	AuthorizeURL string `json:"authorize-url,omitempty" yaml:"authorize-url"`
	// ClientID is the client id for the user
	ClientID string `json:"client-id,omitempty" yaml:"client-id"`
	// ClientSecret is the client secret used for authentication
	ClientSecret string `json:"client-secret,omitempty" yaml:"client-secret"`
	// IDToken the identity token
	IDToken string `json:"id-token,omitempty" yaml:"id-token"`
	// RefreshToken is refresh token for the user
	RefreshToken string `json:"refresh-token,omitempty" yaml:"refresh-token"`
	// TokenURL is the endpoint for tokens
	TokenURL string `json:"token-url,omitempty" yaml:"token-url"`
}

// Profile links endpoint and a credential together
type Profile struct {
	// AuthInfo is the credentials to use
	AuthInfo string `json:"user,omitempty" yaml:"user"`
	// Server is a reference to the server config
	Server string `json:"server,omitempty" yaml:"server"`
	// Team is the default team for this profile
	Team string `json:"team,omitempty" yaml:"team"`
}

// Server defines an endpoint for the api server
type Server struct {
	// Endpoint the url for the api endpoint of kore
	Endpoint string `json:"server,omitempty" yaml:"server"`
}

// RefreshResponse is the response returned when a token is refreshed
type RefreshResponse struct {
	// AccessToken is the access token provided
	AccessToken string `json:"access_token,omitempty" yaml:"access_token"`
	// IDToken string is the identity token
	IDToken string `json:"id_token,omitempty" yaml:"id_token"`
}
