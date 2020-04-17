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

var (
	// AllMethods contains all http methods
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

// Verifier is the interface to a verifier
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Verifier
type Verifier interface {
	// Admit handles the inject anything in the request
	Admit(*http.Request) (bool, error)
}

// Config is the configuration for the service
type Config struct {
	// IDPClientID is the client issuer
	IDPClientID string `json:"idp_client_id,omitempty"`
	// IDPServerURL is the openid server url
	IDPServerURL string `json:"idp_server_url,omitempty"`
	// Listen is the interface to listen on
	Listen string `json:"listen,omitempty"`
	// EnableProxyProtocol indicates we should use proxy protocol
	EnableProxyProtocol bool `json:"enable_proxy_protocol,omitempty"`
	// TLSCaAuthority is a caroot used when verifying the upstream idp
	TLSCaAuthority string `json:"tls_ca_authority,omitempty"`
	// TLSCert is the certificate to serve
	TLSCert string `json:"tls_cert,omitempty"`
	// TLSKey is the private key for the above
	TLSKey string `json:"tls_key,omitempty"`
	// SigningCA is used when not using the IDP server url
	SigningCA string `json:"signing_ca,omitempty"`
	// IDPUserClaims is a collection of claims to extract the user idenity
	IDPUserClaims []string `json:"idp_user_claims,omitempty"`
	// IDPGroupClaims is a colletion of claims to extract the group
	IDPGroupClaims []string `json:"idp_group_claims,omitempty"`
	// MetricsListen is the interface for metrics to render
	MetricsListen string `json:"metrics_listen,omitempty"`
	// UpstreamAuthorizationToken is the upstream authentication token to use
	UpstreamAuthorizationToken string `json:"upstream_authorization_token,omitempty"`
	// UpstreamURL is the endpoint to forward requests
	UpstreamURL string `json:"upstream_url,omitempty"`
	// AllowedIPs contains the allowed IP address ranges which are allowed to connect to the proxy
	AllowedIPs []string `json:"allowed_ips,omitempty"`
}

// Interface is the contract to the proxy
type Interface interface {
	// Run start the proxy up
	Run(context.Context) error
	// Stop calls a halt to the proxy
	Stop() error
	// Addr returns with the server address
	Addr() string
}
