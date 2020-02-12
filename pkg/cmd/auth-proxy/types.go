/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package authproxy

import (
	"context"
	"net/http"
)

// Key is the id token context key
type Key struct{}

// Config is the configuration for the service
type Config struct {
	// ClientID is the client issuer
	ClientID string `json:"client_id,omitempty"`
	// DiscoveryURL is the openid discovery url
	DiscoveryURL string `json:"discovery_url,omitempty"`
	// Listen is the interface to listen on
	Listen string `json:"listen,omitempty"`
	// TLSCaAuthority is a caroot used when verifying the upstream idp
	TLSCaAuthority string `json:"tls_ca_authority,omitempty"`
	// TLSCert is the certificate to serve
	TLSCert string `json:"tls_cert,omitempty"`
	// TLSKey is the private key for the above
	TLSKey string `json:"tls_key,omitempty"`
	// SigningCA is used when not using the discovery url
	SigningCA string `json:"signing_ca,omitempty"`
	// UserClaims is a collection of claims to extract the user idenity
	UserClaims []string `json:"user_claims,omitempty"`
	// GroupClaims is a colletion of claims to extract the group
	GroupClaims []string `json:"group_claims,omitempty"`
	// UpstreamAuthorizationToken is the upstream authentication token to use
	UpstreamAuthorizationToken string `json:"upstream_authorization_token,omitempty"`
	// UpstreamURL is the endpoint to forward requests
	UpstreamURL string `json:"upstream_url,omitempty"`
}

// Interface is the contract to the proxy
type Interface interface {
	// Run start the proxy up
	Run(context.Context) error
	// Stop calls a halt to the proxy
	Stop() error
}

var (
	AllMethods = []string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}
)
