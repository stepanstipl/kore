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

package client

import (
	"context"
	"io"
	"net/http"
)

// Interface is the api client interface
type Interface interface {
	// HTTPClient sets the http client
	HTTPClient(*http.Client) Interface
	// Request creates a request instance
	Request() RestInterface
	// CurrentProfile returns the current profile
	CurrentProfile() string
	// OverrideProfile allows you set the selected profile
	OverrideProfile(string) Interface
}

// Plugin provides an interface for plugins
type Plugin interface {
	// WrapTransport provides the entrypoint for the authentication handle
	WrapTransport(http.RoundTripper) http.RoundTripper
}

// RestInterface provides the rest interface
type RestInterface interface {
	// Body returns the body if any
	Body() io.Reader
	// Context sets the request context
	Context(context.Context) RestInterface
	// Delete performs a delete
	Delete() RestInterface
	// Do returns the response and error
	Do() (RestInterface, error)
	// Endpoint defines the endpoint to use
	Endpoint(string) RestInterface
	// Exists checks if the resource exists
	Exists() (bool, error)
	// Error returns the error if any
	Error() error
	// Get performs a get request
	Get() RestInterface
	// GetPayload returns the payload for inspection
	GetPayload() interface{}
	// HasParamater checks if the parameter is set
	HasParamater(string) (string, bool)
	// Name sets the resource name
	Name(string) RestInterface
	// Resource set the resource kind in the request
	Resource(string) RestInterface
	// Parameters defines a list of parameters for the request
	Parameters(...ParameterFunc) RestInterface
	// Payload set the payload of the request
	Payload(interface{}) RestInterface
	// Result set the object which we should decode into
	Result(interface{}) RestInterface
	// SubResource adds a subresource to the operation
	SubResource(string) RestInterface
	// Team set the team
	Team(string) RestInterface
	// Update performs an put request
	Update() RestInterface
}

// ParameterFunc defines a method for a parameter type
type ParameterFunc func() (Parameter, error)

// Parameter is a param to the raw endpoint
type Parameter struct {
	// IsPath indicates if it's a path or query parameter
	IsPath bool
	// Name is the name of the parameter
	Name string
	// Value is the value of the parameter
	Value string
}
