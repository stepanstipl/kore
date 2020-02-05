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

package apiserver

import (
	"context"

	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime"

	restful "github.com/emicklei/go-restful"
)

const (
	// APIVersion is the hub api server we are serving
	APIVersion = "/api/v1alpha1"
)

// Interface is the interface to the api server
type Interface interface {
	// Run starts the api up
	Run(context.Context) error
	// Stop indicates to want to stop the api
	Stop(context.Context) error
	// BaseURI return the base uri
	BaseURI() string
}

// Config is the configuration for the api server
type Config struct {
	// EnableDex indicates if the idp endpoints should be enabled
	EnableDex bool `json:"enable-dex,omitempty"`
	// Listen is the interface the api should bind on
	Listen string `json:"listen,omitempty"`
	// MetaStoreURL is the host url for the metadata store
	MetaStoreURL string `json:"meta-store-url,omitempty"`
	// MetricsPort is the port the metrics http server should be served
	MetricsPort int `json:"metrics-port,omitempty"`
	// PublicURL is the public url for the api
	PublicURL string `json:"public-url,omitempty"`
	// SwaggerUIPath is the path to the swagger-ui assets
	SwaggerUIPath string `json:"swagger-ui-path,omitempty"`
	// TLSCert is the path the tls certificate
	TLSCert string `json:"tls-cert,omitempty"`
	// TLSKey is the path to the tls private key
	TLSKey string `json:"tls-key,omitempty"`
}

// DefaultHandler implements a default handler
type DefaultHandler struct{}

// EnableAuthentication defaults to yes we do
func (d DefaultHandler) EnableAuthentication() bool {
	return true
}

// EnableLogging defaults to true
func (d DefaultHandler) EnableLogging() bool {
	return true
}

// EnableAdminsOnly requires the user is a member of the admin group
func (d DefaultHandler) EnableAdminsOnly() bool {
	return false
}

// AuthorizationResponse contains the result of a authorization request
type AuthorizationResponse struct {
	// AuthorizationURL is the endpoint for identity provider
	AuthorizationURL string `json:"authorization_url,omitempty"`
	// ClientID is the client id of the login
	ClientID string `json:"client_id,omitempty"`
	// ClientSecret is used for refreshing
	ClientSecret string `json:"client_secret,omitempty"`
	// AccessToken is the access token provided
	AccessToken string `json:"access_token,omitempty"`
	// RefreshToken is a potential refresh token
	RefreshToken string `json:"refresh_token,omitempty"`
	// IDToken string is the identity token
	IDToken string `json:"id_token,omitempty"`
	// TokenEndpointURL is the token endpoint
	TokenEndpointURL string `json:"token_endpoint_url,omitempty"`
}

// Resource defines a resource handler
type Resource interface {
	// Name is the name of resource
	Name() string
	// Kind returns the resource type
	Kind() runtime.Object
	// Delete is called to remove a resource
	Delete(context.Context, string) error
	// Get returns a runtime object
	Get(context.Context, string) (runtime.Object, error)
	// List is called to return a list of the resources
	List(context.Context) ([]runtime.Object, error)
	// Update is the update handler
	Update(context.Context, runtime.Object) error
}

// Handler is the contract to a resource handler
type Handler interface {
	// EnableAdminsOnly requires the user is a member of the admin group
	EnableAdminsOnly() bool
	// EnableAuthentication indicates if the webservice requires authentication
	EnableAuthentication() bool
	// EnableLogging switches of logging for the service
	EnableLogging() bool
	// Name returns the name of the api handler
	Name() string
	// Register is called to allow the handler to register
	Register(hub.Interface, utils.PathBuilder) (*restful.WebService, error)
}
