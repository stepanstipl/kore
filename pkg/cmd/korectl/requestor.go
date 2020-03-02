/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ghodss/yaml"
	"github.com/savaki/jq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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

// Requestor is responsible for calling out to the API
type Requestor struct {
	// config is the cli configuration
	config *Config
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
	// response is the decoded json
	response map[string]interface{}
	// body is the content read bac
	body *bytes.Buffer
	// payload
	payload *bytes.Buffer
	// injections - oh my god - this is what happens when you write things fast
	injections map[string]string
	// if set it will be encoded as JSON as the payload
	runtimeObj runtime.Object
	// responseHandler can be used to register a response handler
	responseHandler func(resp *http.Response) error
}

// NewRequest creates and returns a requestor
func NewRequest() *Requestor {
	return &Requestor{
		params:      make(map[string]bool),
		queryParams: make(map[string]string),
		response:    make(map[string]interface{}),
		injections:  make(map[string]string),
		payload:     nil,
	}
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

	resp, err := c.makeRequest(http.MethodGet, url)
	if err != nil {
		return err
	}

	if err := c.handleResponse(resp); err != nil {
		return err
	}

	return c.parseResponse()
}

// Edit is responsible for performing the request
func (c *Requestor) Edit() error {
	return nil
}

// Update is responsible for performing the request
func (c *Requestor) Update() error {
	url, err := c.makeURI()
	if err != nil {
		return err
	}
	resp, err := c.makeRequest(http.MethodPut, url)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// Delete is responsible for performing the request
func (c *Requestor) Delete() error {
	url, err := c.makeURI()
	if err != nil {
		return err
	}

	resp, err := c.makeRequest(http.MethodDelete, url)
	if err != nil {
		return err
	}

	return c.handleResponse(resp)
}

// parseResponse
func (c *Requestor) parseResponse() error {
	if c.cliCtx.IsSet("output") {
		format := c.cliCtx.String("output")
		switch format {
		case "json":
			return json.NewEncoder(os.Stdout).Encode(c.response)
		case "yaml":
			out, err := yaml.Marshal(c.response)
			fmt.Fprintf(os.Stdout, "%s", out)
			return err
		default:
			return errors.New("unsupported output type")
		}
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 20, 0, ' ', 10)

	columns := strings.Join(c.render, "\t")
	fmt.Fprintf(w, "%s\n", columns)

	islist := c.response["items"] != nil
	switch islist {
	case true:
		items, found := c.response["items"].([]interface{})
		if !found {
			return errors.New("invalid response list, no items found")
		}
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
	default:
		values, err := c.makeValues(c.Body(), c.paths)
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

// handleResponse is used to wrap common errors
func (c Requestor) handleResponse(resp *http.Response) error {
	if c.responseHandler != nil {
		return c.responseHandler(resp)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// @step: does the error contain custom error?
	if c.response["code"] != nil && c.response["message"] != "" {
		fmt.Printf("[error] [%d] %s\n", int(c.response["code"].(float64)), c.response["message"])
		os.Exit(1)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		fmt.Println("[error] authorization required, please use the 'authorize' command")
	case http.StatusNotFound:
		fmt.Println("[error] from server (notfound): resource not found")
	case http.StatusForbidden:
		fmt.Println("[error] request has been denied, check credentials")
	case http.StatusBadRequest:
		fmt.Println("[error] api responded with invalid request")
	default:
		fmt.Printf("invalid response: %d from api server", resp.StatusCode)
	}
	os.Exit(1)

	return nil
}

// makeRequest creates and returns a http request
func (c *Requestor) makeRequest(method, url string) (*http.Response, error) {
	var req *http.Request
	var err error

	if c.runtimeObj != nil {
		encoded, err := json.Marshal(c.runtimeObj)
		if err != nil {
			return nil, err
		}
		c.payload = bytes.NewBuffer(encoded)
	}

	if c.payload == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, c.payload)
	}
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("kind or resource does not exist")
	}

	if resp.Body != nil {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if len(content) > 0 {
			if err := json.NewDecoder(bytes.NewReader(content)).Decode(&c.response); err != nil {
				return nil, err
			}
		}
		c.body = bytes.NewBuffer(content)
	}

	return resp, nil
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

func (c *Requestor) WithRuntimeObject(obj runtime.Object) *Requestor {
	c.runtimeObj = obj
	return c
}

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

func (c *Requestor) WithResponseHandler(f func(resp *http.Response) error) *Requestor {
	c.responseHandler = f
	return c
}

// WithConfig adds the configuration
func (c *Requestor) WithConfig(config *Config) *Requestor {
	c.config = config

	return c
}

func (c *Requestor) Body() io.Reader {
	return c.body
}
