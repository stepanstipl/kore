/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
	// Server is the api server url
	Server string `json:"server,omitempty" yaml:"server"`
	// Credentials are the credentials for the api server
	Credentials Identity `json:"credentials,omitempty" yaml:"credentials"`
	// TokenURL is the endpoint for tokens
	TokenURL string `json:"token-url,omitempty" yaml:"token_url"`
	// AuthorizeURL is the endpoint for the authorize
	AuthorizeURL string `json:"authorize-url,omitempty" yaml:"authorize_url"`
}

type AuthorizationResponse struct {
	// AuthorizationURL is the endpoint for identity provider
	AuthorizationURL string `json:"authorization_url,omitempty" yaml:"authorization_url"`
	// ClientID is the client id of the login
	ClientID string `json:"client_id,omitempty" yaml:"client_id"`
	// ClientSecret is used for refreshing
	ClientSecret string `json:"client_secret,omitempty" yaml:"client_secret"`
	// AccessToken is the access token provided
	AccessToken string `json:"access_token,omitempty" yaml:"access_token"`
	// RefreshToken is a potential refresh token
	RefreshToken string `json:"refresh_token,omitempty" yaml:"refresh_token"`
	// IDToken string is the identity token
	IDToken string `json:"id_token,omitempty" yaml:"id_token"`
	// TokenEndpointURL is the token endpoint
	TokenEndpointURL string `json:"token_endpoint_url,omitempty" yaml:"token_endpoint_url"`
}

// Identity is the identity within the hub
type Identity struct {
	// AccessToken is the access token retrieved from hub
	AccessToken string `json:"access-token,omitempty" yaml:"access_token"`
	// ClientID is the client id for the user
	ClientID string `json:"client-id,omitempty" yaml:"client_id"`
	// ClientSecret is the client secret used for authentication
	ClientSecret string `json:"client-secret,omitempty" yaml:"client_secret"`
	// IDToken the identity token
	IDToken string `json:"id_token,omitempty" yaml:"id_token"`
	// RefreshToken is refresh token for the user
	RefreshToken string `json:"refresh-token,omitempty" yaml:"refresh_token"`
}
