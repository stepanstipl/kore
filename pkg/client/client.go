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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/client/config"
	cerrors "github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/utils/validation"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
)

// apiClient implements the api and raw client
type apiClient struct {
	// body is the response body
	body *bytes.Buffer
	// ctx is the context for the client
	ctx context.Context
	// cfg is the client configuration
	cfg *config.Config
	// ferror is used to handle errors in the method chain
	ferror error
	// endpoint is the raw endpoint template to use
	endpoint string
	// hc is the http client to use
	hc *http.Client
	// profile is the name of the profile to use
	profile string
	// parameters hold the parameters for the request
	parameters map[string]string
	// payload is the outbound payload
	payload interface{}
	// result is what we decode into
	result interface{}
	// queryparams are a collection of query parameters
	queryparams url.Values
}

var (
	// ErrAuthenticationRequired requires authentication
	ErrAuthenticationRequired = apiserver.Error{
		Code:    http.StatusUnauthorized,
		Message: "authentication required",
	}
)

// cc provides a wrapper around th config
type cc struct {
	cfg     *config.Config
	hc      *http.Client
	profile string
}

// New creates and returns an API client
func New(c *config.Config) (Interface, error) {
	if c == nil {
		return nil, errors.New("no client configuration")
	}

	return &cc{cfg: c, hc: DefaultHTTPClient, profile: c.CurrentProfile}, nil
}

// HTTPClient sets the http client
func (c *cc) HTTPClient(hc *http.Client) Interface {
	c.hc = hc

	return c
}

// OverrideProfile sets the default profile to use
func (c *cc) OverrideProfile(name string) Interface {
	c.profile = name

	return c
}

// CurrentProfile returns the current profile
func (c *cc) CurrentProfile() string {
	return c.profile
}

// Request creates a request instance
func (c *cc) Request() RestInterface {
	return &apiClient{
		cfg:        c.cfg,
		hc:         c.hc,
		parameters: make(map[string]string),
		profile:    c.profile,
	}
}

func (a *apiClient) Profile() string {
	return a.profile
}

// HandleRequest is responsible for handling the request chain
func (a *apiClient) HandleRequest(method string) RestInterface {
	err := func() error {
		// @step: check if we had any errors in the method chain
		if a.ferror != nil {
			return a.ferror
		}

		// @step: check we have the endpoint
		profile, found := a.cfg.Profiles[a.Profile()]
		if !found {
			return cerrors.ErrMissingProfile
		}
		server, found := a.cfg.Servers[profile.Server]
		if !found {
			return cerrors.NewProfileInvalidError("missing profile server", a.Profile())
		}
		endpoint := server.Endpoint

		if endpoint == "" {
			return cerrors.NewProfileInvalidError("missing endpoint", a.Profile())
		}

		// @step: we generate the uri from the parameter
		uri, err := a.MakeRequestURI()
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"endpoint": endpoint,
			"method":   method,
			"uri":      uri,
		}).Debug("making request to kore api")

		// @step: we generate the fully qualifies url
		ep := fmt.Sprintf("%s/%s", endpoint, uri)

		// @step: we make the request
		now := time.Now()
		resp, err := a.MakeRequest(method, ep)
		if err != nil {
			return err
		}
		log.WithField("time", time.Since(now).String()).Debug("processing time of the request was")

		return a.HandleResponse(resp)
	}()
	if err != nil {
		a.ferror = err
	}

	return a
}

// MakeRequestURI is responsible for generating the request URI
func (a *apiClient) MakeRequestURI() (string, error) {
	if a.endpoint == "" {
		return a.MakeDefaultURL()
	}

	return a.MakeEndpointURL()
}

// MakeDefaultURL generates a URL from the /teams/<name>/resource/<name> format
func (a *apiClient) MakeDefaultURL() (string, error) {
	var paths []string
	// @logic: we simply iterate in a known order of things i.e. team, resource
	// and name check if the parameter is there and if so append

	if value, found := a.HasParameter("team"); found && value != "" {
		paths = append(paths, []string{"teams", value}...)
	}
	for _, x := range []string{"resource", "name", "subresource"} {
		if value, found := a.HasParameter(x); found {
			paths = append(paths, value)
		}
	}

	baseuri := strings.TrimPrefix(apiserver.APIVersion, "/")

	// @step: we add the path elements and the queries together
	uri := strings.Join(append([]string{baseuri}, paths...), "/")
	if len(a.queryparams) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, a.queryparams.Encode())
	}

	return uri, nil
}

// MakeEndpointURL is responsible for generating the url from a template
func (a *apiClient) MakeEndpointURL() (string, error) {
	uri := strings.TrimPrefix(apiserver.APIVersion, "/") + a.endpoint

	// @step: we add the path elements and the queries together
	for param, value := range a.parameters {
		uri = strings.ReplaceAll(uri, "{"+param+"}", value)
	}

	// @step: add the query params if any to the url
	if len(a.queryparams) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, a.queryparams.Encode())
	}

	return uri, nil
}

// MakeRequest is responsible for preparing and handling the http request
func (a *apiClient) MakeRequest(method, url string) (*http.Response, error) {
	// @step: do we have any thing to encode?
	payload, err := a.MakePayload()
	if err != nil {
		return nil, err
	}

	// @step: construct the http request
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Client-Version", version.Release)

	// @step: add the authentication from profile
	auth := a.cfg.AuthInfos[a.Profile()]
	switch {
	case auth == nil:
		return nil, cerrors.NewProfileInvalidError("missing authenication profile", a.Profile())
	case auth.OIDC != nil:
		request.Header.Set("Authorization", "Bearer "+auth.OIDC.IDToken)
	case auth.Token != nil:
		request.Header.Set("Authorization", "Bearer "+*auth.Token)
	case auth.BasicAuth != nil:
		request.SetBasicAuth(auth.BasicAuth.Username, auth.BasicAuth.Password)
	}

	return a.hc.Do(request)
}

// HandleResponse is responsible for handling the http response from api
func (a *apiClient) HandleResponse(resp *http.Response) error {
	// @step: if everything is ok, check for a response and return
	code := resp.StatusCode

	log.WithField("code", code).Debug("received the following http response code")

	if code >= http.StatusOK && code <= 299 {
		if err := a.MakeResult(resp, a.result); err != nil {
			return err
		}

		return nil
	}

	// @step: we have encountered an error we need read in a APIError or create one
	apiError := &apiserver.Error{}
	apiError.Code = resp.StatusCode
	apiError.Verb = resp.Request.Method
	apiError.URI = resp.Request.RequestURI

	a.decodeError(resp, apiError)

	if apiError.Message == "" {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			apiError.Message = "authorization required"
		case http.StatusNotFound:
			apiError.Message = "kind or resource does not exist"
		case http.StatusForbidden:
			apiError.Message = "request denied, check permissions"
		case http.StatusBadRequest:
			apiError.Message = "api responded with invalid request"
		default:
			apiError.Message = "invalid response from api server"
		}
	}

	return apiError
}

func (a *apiClient) decodeError(resp *http.Response, apiError *apiserver.Error) {
	if resp.Body != nil {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			vError := &validation.Error{}
			if err := a.MakeResult(resp, vError); err != nil {
				log.WithError(err).Debug("response can not be decoded into a validation error")
				return
			}
			apiError.Message = vError.Error()
			return
		case http.StatusConflict:
			err := &validation.ErrDependencyViolation{}
			if err := a.MakeResult(resp, err); err != nil {
				log.WithError(err).Debug("response can not be decoded into a dependency validation error")
				return
			}
			apiError.Message = err.Error()
			return
		}

		err := a.MakeResult(resp, apiError)
		if err != nil {
			log.WithError(err).Debug("response can not be decoded")
		}
	}
}

// MakeResult is responsible for reading the resulting payload
func (a *apiClient) MakeResult(resp *http.Response, data interface{}) error {
	a.body = &bytes.Buffer{}

	if resp.Body == nil {
		log.Trace("request to api had no response body present")

		return nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	a.body.Write(content)

	log.WithField("body", a.body.String()).Trace("we received the following response from api")

	if data == nil {
		log.Trace("no result has been set to save the payload")

		return nil
	}

	return json.NewDecoder(a.Body()).Decode(data)
}

// MakePayload is responsible for encoding the payload if any
func (a *apiClient) MakePayload() (io.Reader, error) {
	if a.payload == nil {
		return nil, nil
	}
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(a.payload); err != nil {
		return nil, err
	}
	log.WithField("payload", b.String()).Trace("using the attached payload for request")

	return b, nil
}

// HasParameter checks if the parameter exists
func (a *apiClient) HasParameter(key string) (string, bool) {
	value, found := a.parameters[key]

	return value, (found && value != "")
}

// Exists check is the resource exists
func (a *apiClient) Exists() (bool, error) {
	if err := a.Get().Error(); err != nil {
		if !IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// Delete performs a delete
func (a *apiClient) Delete() RestInterface {
	return a.HandleRequest(http.MethodDelete)
}

// Get performs a get request
func (a *apiClient) Get() RestInterface {
	return a.HandleRequest(http.MethodGet)
}

// Update performs an put request
func (a *apiClient) Update() RestInterface {
	return a.HandleRequest(http.MethodPut)
}

// Parameters defines a list of parameters for the request
func (a *apiClient) Parameters(params ...ParameterFunc) RestInterface {
	for _, fn := range params {
		param, err := fn()
		if err != nil {
			a.ferror = err

			continue
		}
		if param.IsPath {
			a.parameters[param.Name] = param.Value
		} else {
			if a.queryparams == nil {
				a.queryparams = url.Values{}
			}
			a.queryparams.Add(param.Name, param.Value)
		}
	}

	return a
}

// SubResource adds a subresource to the operation
func (a *apiClient) SubResource(v string) RestInterface {
	a.parameters["subresource"] = v

	return a
}

// Payload set the payload of the request
func (a *apiClient) Payload(v interface{}) RestInterface {
	a.payload = v

	return a
}

// Result set the object which we should decode into
func (a *apiClient) Result(v interface{}) RestInterface {
	a.result = v

	return a
}

// InjectParam is responsible for injecting the parameter
func (a *apiClient) InjectParam(key, value string) {
	if value == "" {
		panic(fmt.Errorf("%q path parameter can not be empty", key))
	}

	a.parameters[key] = value
}

// Name sets the resource name
func (a *apiClient) Name(v string) RestInterface {
	a.InjectParam("name", v)

	return a
}

// Resource set the resource kind in the request
func (a *apiClient) Resource(v string) RestInterface {
	a.InjectParam("resource", v)

	return a
}

// Team set the team
func (a *apiClient) Team(v string) RestInterface {
	a.parameters["team"] = v

	return a
}

// Endpoint defines the endpoint to use
func (a *apiClient) Endpoint(v string) RestInterface {
	a.endpoint = v

	return a
}

// Context sets the request context
func (a *apiClient) Context(ctx context.Context) RestInterface {
	a.ctx = ctx

	return a
}

// Error return any error and resets post
func (a *apiClient) Error() error {
	// we need to reset the error
	defer func() {
		a.ferror = nil
	}()

	return a.ferror
}

// Do returns both the response and error
func (a *apiClient) Do() (RestInterface, error) {
	return a, a.ferror
}

// GetPayload returns the payload for inspection
func (a *apiClient) GetPayload() interface{} {
	return a.payload
}

// Body returns the body if any
func (a *apiClient) Body() io.Reader {
	return strings.NewReader(a.body.String())
}

// HasParamater checks if the parameter is set
func (a *apiClient) HasParamater(key string) (string, bool) {
	value, found := a.parameters[key]

	return value, found
}
