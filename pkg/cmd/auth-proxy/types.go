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
	"time"
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

// Config is the configuration for the service
type Config struct {
	// AllowedIPs is the collection of networks permitted access
	AllowedIPs []string `json:"allowed-ips,omitempty"`
	// EnableProxyProtocol indicates we should use proxy protocol
	EnableProxyProtocol bool `json:"enable-proxy-protocol,omitempty"`
	// FlushInterval is flush interval for the proxy
	FlushInterval time.Duration `json:"flush-interval,omitempty"`
	// MetricsListen is the interface for metrics to render
	MetricsListen string `json:"metrics-listen,omitempty"`
	// Listen is the interface to listen on
	Listen string `json:"listen,omitempty"`
	// TLSCaAuthority is a ca used when verifying the upstream idp
	TLSCaAuthority string `json:"tls-ca-authority,omitempty"`
	// TLSCert is the certificate to serve
	TLSCert string `json:"tls-cert,omitempty"`
	// TLSKey is the private key for the above
	TLSKey string `json:"tls-key,omitempty"`
	// Token is the kubernetes token
	Token string `json:"token"`
	// UpstreamURL is the endpoint to forward requests
	UpstreamURL string `json:"upstream-url,omitempty"`
	// Verifiers is a collection of verifiers to switch on
	Verifiers []string `json:"verifiers,omitempty"`
}

// Interface is the contract to the proxy
type Interface interface {
	// Run start the proxy up
	Run(context.Context) error
	// Stop calls a halt to the proxy
	Stop() error
}
