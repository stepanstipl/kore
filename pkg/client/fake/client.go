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

package fake

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/appvia/kore/pkg/client"
)

// Request creates a request instance
func (f *fake) Request() client.RestInterface {
	return f
}

// HTTPClient sets the http client
func (f *fake) HTTPClient(*http.Client) client.Interface {
	return f
}

// OverrideProfile does nothing here
func (f *fake) OverrideProfile(string) client.Interface {
	return f
}

// CurrentProfile returns the current profile
func (f *fake) CurrentProfile() string {
	return ""
}

// Body returns the body if any
func (f *fake) Body() io.Reader {
	b := &bytes.Buffer{}

	if f.result != nil {
		_ = json.NewEncoder(b).Encode(f.result)
	}

	return b
}

// Context sets the request context
func (f *fake) Context(context.Context) client.RestInterface {
	return f
}

// Delete performs a delete
func (f *fake) Delete() client.RestInterface {
	return f
}

// SubResource is an operation under the resource
func (f *fake) SubResource(string) client.RestInterface {
	return f
}

// Do returns the response and error
func (f *fake) Do() (client.RestInterface, error) {
	return f, nil
}

// Endpoint defines the endpoint to use
func (f *fake) Endpoint(string) client.RestInterface {
	return f
}

// Exists checks if the resource exists
func (f *fake) Exists() (bool, error) {
	return false, nil
}

// Error returns the error if any
func (f *fake) Error() error {
	return nil
}

// Get performs a get request
func (f *fake) Get() client.RestInterface {
	return f
}

// GetPayload returns the payload for inspection
func (f *fake) GetPayload() interface{} {
	return f.payload
}

// HasParamater checks if the parameter is set
func (f *fake) HasParamater(string) (string, bool) {
	return "", false
}

// Name sets the resource name
func (f *fake) Name(string) client.RestInterface {
	return f
}

// Resource set the resource kind in the request
func (f *fake) Resource(string) client.RestInterface {
	return f
}

// Parameters defines a list of parameters for the request
func (f *fake) Parameters(...client.ParameterFunc) client.RestInterface {
	return f
}

// Payload set the payload of the request
func (f *fake) Payload(interface{}) client.RestInterface {
	return f
}

// Result set the object which we should decode into
func (f *fake) Result(interface{}) client.RestInterface {
	return f
}

// Team set the team
func (f *fake) Team(string) client.RestInterface {
	return f
}

// Update performs an put request
func (f *fake) Update() client.RestInterface {
	return f
}
