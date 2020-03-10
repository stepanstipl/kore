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

package authproxy

import (
	"context"
	"net/http"
)

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
	// MetricsListen is the interface for metrics to render
	MetricsListen string `json:"metrics_listen,omitempty"`
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
