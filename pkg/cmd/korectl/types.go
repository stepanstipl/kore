/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
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

package korectl

import resty "gopkg.in/resty.v1"

var (
	// DefaultHome is the home directory for the korectl
	DefaultHome = "${HOME}/.korectl"
	// HubConfig is the configuration file for the api
	HubConfig = DefaultHome + "/config"
)

var (
	hp = resty.New()
)

// Config is the configuration for the api
type Config struct {
	// CurrentContext is the context in use at the moment
	CurrentContext string `json:"current-context,omitempty" yaml:"current-context"`
	// Contexts is a collection of contexts
	Contexts map[string]*Context `json:"contexts,omitempty" yaml:"contexts"`
	// Servers is a collection of api endpoints
	Servers map[string]*Server `json:"servers,omitempty" yaml:"servers"`
	// AuthInfos is a collection of credentials
	AuthInfos map[string]*AuthInfo `json:"users,omitempty" yaml:"users"`
}

// AuthInfo defines a credential to the api endpoint
type AuthInfo struct {
	// Token is a static token to use
	Token *string `json:"token,omitempty" yaml:"token"`
	// OIDC is credentials from an oauth2 provider
	OIDC *OIDC `json:"oidc" yaml:"oidc"`
}

// OIDC is the identity within the kore
type OIDC struct {
	// AccessToken is the access token retrieved from kore
	AccessToken string `json:"access-token,omitempty" yaml:"access-token"`
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
	// AuthorizeURL is the endpoint for the authorize
	AuthorizeURL string `json:"authorize-url,omitempty" yaml:"authorize-url"`
}

// Context links endpoint and a credential together
type Context struct {
	// Server is a reference to the server config
	Server string `json:"server,omitempty" yaml:"server"`
	// AuthInfo is the credentials to use
	AuthInfo string `json:"user,omitempty" yaml:"user"`
}

// Server defines an endpoint for the api servr
type Server struct {
	// Endpoint is the server url
	Endpoint string `json:"server,omitempty" yaml:"server"`
}

type AuthorizationResponse struct {
	// AuthorizationURL is the endpoint for identity provider
	AuthorizationURL string `json:"authorization-url,omitempty" yaml:"authorization-url"`
	// ClientID is the client id of the login
	ClientID string `json:"client-id,omitempty" yaml:"client-id"`
	// ClientSecret is used for refreshing
	ClientSecret string `json:"client-secret,omitempty" yaml:"client-secret"`
	// AccessToken is the access token provided
	AccessToken string `json:"access-token,omitempty" yaml:"access-token"`
	// RefreshToken is a potential refresh token
	RefreshToken string `json:"refresh-token,omitempty" yaml:"refresh-token"`
	// IDToken string is the identity token
	IDToken string `json:"id-token,omitempty" yaml:"id-token"`
	// TokenEndpointURL is the token endpoint
	TokenEndpointURL string `json:"token-endpoint-url,omitempty" yaml:"token-endpoint-url"`
}
