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

package korectl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/appvia/kore/pkg/kore/validation"
	"github.com/appvia/kore/pkg/utils"

	"github.com/ghodss/yaml"
	"github.com/savaki/jq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	hc = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
)

type RequestError struct {
	err        error
	statusCode int
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("[%d] %s", r.statusCode, r.err.Error())
}

func (r *RequestError) StatusCode() int {
	return r.statusCode
}

// Requestor is responsible for calling out to the API
type Requestor struct {
	// config is the cli configuration
	config *Config
	// configError tracks errors during the configuration process
	configError error
	// cliCtx is the cli context
	cliCtx *cli.Context
	// endpoint is the uri endpoint to call
	endpoint string
	// params are the parameters off the cli to substitute
	params map[string]bool
	// queryParams are query params
	queryParams map[string]string
	// render is the render paths
	render []string
	// paths are the paths for render
	paths []string
	// payload
	payload *bytes.Buffer
	// injections - oh my god - this is what happens when you write things fast
	injections map[string]string
	// if set it will be encoded as JSON as the payload
	runtimeObj interface{}
}

// NewRequest creates and returns a requestor
func NewRequest() *Requestor {
	return &Requestor{
		params:      make(map[string]bool),
		queryParams: make(map[string]string),
		injections:  make(map[string]string),
		payload:     nil,
	}
}

func NewRequestForResource(config *Config, ctx *cli.Context) (*Requestor, resourceConfig, error) {
	resConfig := getResourceConfig(ctx.Args().First())

	req := NewRequest().
		WithConfig(config).
		WithContext(ctx).
		PathParameter("resource", true).
		WithInject("resource", resConfig.Name)

	var endpoint string

	if ctx.IsSet("team") {
		endpoint = "/teams/{team}/{resource}"
		if !resConfig.IsTeam {
			return nil, resourceConfig{}, errors.New("--team parameter is not allowed for this resource")
		}
		req.PathParameter("team", true)
	} else {
		endpoint = "/{resource}"
		if !resConfig.IsGlobal {
			return nil, resourceConfig{}, errTeamParameterMissing
		}
	}

	if ctx.NArg() == 2 {
		endpoint = endpoint + "/{name}"
		req.PathParameter("name", true)
		req.WithInject("name", ctx.Args().Get(1))
	}

	req.WithEndpoint(endpoint)

	return req, resConfig, nil
}

// Column sets a column option
func Column(name, path string) string {
	return fmt.Sprintf("%s/%s", name, path)
}

// Get is responsible for performing the request
func (c *Requestor) Get() error {
	url, err := c.makeURI()
	if err != nil {
		return err
	}

	responseHandler := c.parseResponse
	if c.runtimeObj != nil {
		responseHandler = c.parseObjectResponse
	}

	return c.doRequest(http.MethodGet, url, responseHandler)
}

// Exists will perform a GET request and will return
//  * true on 200 response
//  * false on 404 response
//  * error on any other response
func (c *Requestor) Exists() (bool, error) {
	url, err := c.makeURI()
	if err != nil {
		return false, err
	}

	if err := c.doRequest(http.MethodGet, url, nil); err != nil {
		if reqErr, ok := err.(*RequestError); ok {
			if reqErr.statusCode == http.StatusNotFound {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

// Update is responsible for performing the request
func (c *Requestor) Update() error {
	url, err := c.makeURI()
	if err != nil {
		return err
	}
	return c.doRequest(http.MethodPut, url, nil)
}

// Delete is responsible for performing the request
func (c *Requestor) Delete() error {
	url, err := c.makeURI()
	if err != nil {
		return err
	}

	return c.doRequest(http.MethodDelete, url, nil)
}

func (c *Requestor) parseObjectResponse(resp *http.Response) error {
	if err := json.NewDecoder(resp.Body).Decode(c.runtimeObj); err != nil {
		return err
	}
	return nil
}

func (c *Requestor) parseResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	content := bytes.NewReader(body)

	var response map[string]interface{}
	if content.Len() > 0 {
		if err := json.NewDecoder(content).Decode(&response); err != nil {
			return err
		}
	}
	if c.cliCtx == nil {
		return nil
	}
	if c.cliCtx.IsSet("output") {
		format := c.cliCtx.String("output")
		switch format {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(response)
		case "yaml":
			out, err := yaml.Marshal(response)
			fmt.Fprintf(os.Stdout, "%s", out)
			return err
		default:
			return errors.New("unsupported output type")
		}
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 10, 4, ' ', 0)

	columns := strings.Join(c.render, "\t")
	fmt.Fprintf(w, "%s\n", columns)

	if _, found := response["items"]; found {
		items, found := response["items"].([]interface{})
		if found {
			for _, x := range items {
				decoded := &bytes.Buffer{}
				if err := json.NewEncoder(decoded).Encode(x); err != nil {
					return err
				}
				values, err := c.makeValues(decoded, c.paths)
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "%s\n", strings.Join(values, "\t"))
			}
		}
	} else {
		_, _ = content.Seek(0, io.SeekStart)
		values, err := c.makeValues(content, c.paths)
		if err != nil {
			return err

		}
		fmt.Fprintf(w, "%s", strings.Join(values, "\t"))
	}

	return w.Flush()
}

// makeValues returns the jsonpath values
func (c *Requestor) makeValues(in io.Reader, paths []string) ([]string, error) {
	var list []string

	if in == nil {
		return []string{}, errors.New("no request body")
	}
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return []string{}, err
	}

	for _, x := range paths {
		op, err := jq.Parse(x)
		if err != nil {
			return list, fmt.Errorf("invalid jsonpath expression: %s, error: %s", x, err)
		}
		v, err := op.Apply(data)
		if err != nil {
			if !strings.Contains(err.Error(), "key not found") {
				return list, fmt.Errorf("failed to apply jsonpath to response body: %s", err)
			}
			v = []byte("Unknown")
		}

		list = append(list, strings.ReplaceAll(string(v), "\"", ""))
	}

	return list, nil
}

// checkResponse is used to check whether the request was successful
func (c Requestor) checkResponse(resp *http.Response) *RequestError {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		// This is a validation error, check if we have a validation error body. If so,
		// that implements Error so we can use it directly.
		var valResponse validation.ErrValidation
		if err := json.NewDecoder(resp.Body).Decode(&valResponse); err == nil {
			return &RequestError{
				statusCode: resp.StatusCode,
				err:        valResponse,
			}
		}
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		// @step: does the error contain custom error?
		if response["code"] != nil && response["message"] != "" {
			return &RequestError{
				statusCode: int(response["code"].(float64)),
				err:        fmt.Errorf("%s", response["message"]),
			}
		}
	}

	var err error
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		err = errors.New("authorization required, please use the 'login' command")
	case http.StatusNotFound:
		err = errors.New("kind or resource does not exist")
	case http.StatusForbidden:
		err = errors.New("request has been denied, check credentials")
	case http.StatusBadRequest:
		err = errors.New("api responded with invalid request")
	default:
		err = fmt.Errorf("invalid response: %d from api server", resp.StatusCode)
	}

	return &RequestError{
		statusCode: resp.StatusCode,
		err:        err,
	}
}

// doRequest makes and handles an HTTP request
func (c *Requestor) doRequest(method, url string, handler func(*http.Response) error) error {
	if c.configError != nil {
		return c.configError
	}

	var req *http.Request
	var err error

	if c.runtimeObj != nil {
		encoded, err := json.Marshal(c.runtimeObj)
		if err != nil {
			return err
		}
		c.payload = bytes.NewBuffer(encoded)
	}

	if c.payload == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, c.payload)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	auth := c.config.GetCurrentAuthInfo()
	if auth.Token != nil {
		req.Header.Set("Authorization", "Bearer "+*auth.Token)
	}
	if auth.OIDC != nil {
		req.Header.Set("Authorization", "Bearer "+auth.OIDC.IDToken)
	}

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	if handler != nil {
		err := handler(resp)
		if err != nil {
			return &RequestError{statusCode: resp.StatusCode, err: err}
		}
		return nil
	} else {
		// Read the response body and discard it - to avoid any surprises
		_, _ = ioutil.ReadAll(resp.Body)

		return nil
	}
}

// makeURI is responsible for generating the uri for the requestor
func (c *Requestor) makeURI() (string, error) {
	uri := c.endpoint

	for name, required := range c.params {
		value, found := c.injections[name]
		if !found {
			found = c.cliCtx.IsSet(name)
			value = fmt.Sprintf("%s", c.cliCtx.Generic(name))
		}
		if !found && required {
			return "", fmt.Errorf("invalid request, option: %s must be set", name)
		}
		token := fmt.Sprintf("{%s}", name)
		if found {
			if !strings.Contains(uri, token) {
				uri = fmt.Sprintf("%s/%s", uri, value)
			} else {
				uri = strings.ReplaceAll(uri, token, value)
			}
		} else {
			if strings.Contains(uri, token) {
				uri = strings.ReplaceAll(uri, token, "")
			}
		}
	}
	uri = strings.TrimSuffix(uri, "/")
	uri = strings.TrimPrefix(uri, "/")

	if len(c.queryParams) <= 0 {
		return fmt.Sprintf("%s/%s", c.config.GetAPI(), uri), nil
	}
	var list []string

	for k, v := range c.queryParams {
		list = append(list, fmt.Sprintf("%s=%v", k, v))
	}
	url := fmt.Sprintf("%s/%s?%s", c.config.GetAPI(), uri, strings.Join(list, "&"))

	log.WithField("url", url).Debug("making request to kore apiserver")

	return url, nil
}

// Render is used to set a render layout
func (c *Requestor) Render(opts ...string) *Requestor {
	for _, x := range opts {
		items := strings.Split(x, "/")
		if len(items) == 2 {
			c.render = append(c.render, items[0])
			c.paths = append(c.paths, items[1])
		}
	}

	return c
}

// QueryParam sets a query parameters
func (c *Requestor) QueryParam(param string, value interface{}) *Requestor {
	c.queryParams[param] = fmt.Sprintf("%s", value)

	return c
}

// PathParameter sets a path parameters
func (c *Requestor) PathParameter(param string, required bool) *Requestor {
	c.params[param] = required

	return c
}

// WithInject
func (c *Requestor) WithInject(name, value string) *Requestor {
	c.injections[name] = value

	return c
}

// WithEndpoint sets the endpoint
func (c *Requestor) WithEndpoint(uri string) *Requestor {
	c.endpoint = uri

	return c
}

// WithContext adds the cli context
func (c *Requestor) WithContext(ctx *cli.Context) *Requestor {
	c.cliCtx = ctx

	return c
}

func (c *Requestor) WithRuntimeObject(obj interface{}) *Requestor {
	c.runtimeObj = obj

	return c
}

// WithPayload sets the payload from a file
func (c *Requestor) WithPayload(name string) *Requestor {
	path := c.cliCtx.String(name)
	if path == "" {
		return c
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("trying to read file: %s", path))
	}
	encoded, err := yaml.YAMLToJSON(content)
	if err != nil {
		panic(fmt.Sprintf("trying to parse the yaml to json: %s", path))
	}
	u := &unstructured.Unstructured{}
	if err := u.UnmarshalJSON(encoded); err != nil {
		panic(fmt.Sprintf("trying to parse the contents of file: %s", err))
	}

	if u.GetName() == "" {
		panic("resource does not have a name")
	}

	c.injections["name"] = u.GetName()
	c.payload = bytes.NewBuffer(encoded)

	return c
}

// HasExpiredToken checks if the token in an oidc config is invalid
func (o *OIDC) HasExpiredToken() bool {
	claims, _ := utils.NewClaimsFromRawToken(o.IDToken)
	exp, found := claims.GetExpiry()
	// The token has no exp claim and is therefore invalid
	if !found {
		panic("trying to validate the openid token: exp claim is not set")
	}
	return exp.Before(time.Now().UTC())
}

// WithConfig adds the configuration
func (c *Requestor) WithConfig(config *Config) *Requestor {
	auth := config.GetCurrentAuthInfo()

	if auth.OIDC != nil {
		if auth.OIDC.HasExpiredToken() {
			log.Debug("token has expired, requesting a new one")

			if err := auth.OIDC.Refresh(); err != nil {
				c.configError = err
			}

			if err := config.Update(); err != nil {
				panic("trying to update the config file")
			}
		}
	}

	c.config = config
	return c
}
