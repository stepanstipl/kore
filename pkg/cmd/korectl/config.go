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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/ghodss/yaml"
)

// GetSwaggerCachedFile returns the location of the swagger cache
func (c *Config) GetSwaggerCachedFile() string {
	return path.Join(os.ExpandEnv(c.GetDirectory()), "cache.json")
}

// GetDirectory returns the korectl home
func (c *Config) GetDirectory() string {
	if os.Getenv("HUB_CLI_HOME") != "" {
		return os.Getenv("HUB_CLI_HOME")
	}

	return DefaultHome
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() error {
	return nil
}

// SetCurrentProfile is used to set the current profile
func (c *Config) SetCurrentProfile(name string) {
	c.CurrentProfile = name
}

// CreateProfile is used to create a profile
func (c *Config) CreateProfile(name, endpoint string) {
	c.AddProfile(name, &Profile{
		Server:   name,
		AuthInfo: name,
	})
	if !c.HasServer(name) {
		c.AddServer(name, &Server{Endpoint: endpoint})
	}
	if !c.HasAuthInfo(name) {
		c.AddAuthInfo(name, &AuthInfo{OIDC: &OIDC{}})
	}
}

// GetCurrentProfile returns the current profile
func (c *Config) GetCurrentProfile() *Profile {
	profile, found := c.Profiles[c.CurrentProfile]
	if !found {
		return &Profile{}
	}

	return profile
}

// GetCurrentServer returns the server in the context
func (c *Config) GetCurrentServer() *Server {
	ct := c.Profiles[c.CurrentProfile]
	if ct == nil {
		return &Server{}
	}
	s := c.Servers[ct.Server]
	if s == nil {
		return &Server{}
	}

	return s
}

// GetCurrentAuthInfo returns the current auth
func (c *Config) GetCurrentAuthInfo() *AuthInfo {
	ct := c.Profiles[c.CurrentProfile]
	if ct == nil {
		return &AuthInfo{}
	}

	a := c.AuthInfos[ct.AuthInfo]

	if a == nil {
		return &AuthInfo{}
	}

	return a
}

// AddProfile adds a the profile to the config
func (c *Config) AddProfile(name string, ctx *Profile) {
	if c.Profiles == nil {
		c.Profiles = make(map[string]*Profile)
	}
	c.Profiles[name] = ctx
}

// AddServer adds a server
func (c *Config) AddServer(name string, server *Server) {
	if c.Servers == nil {
		c.Servers = make(map[string]*Server)
	}
	c.Servers[name] = server
}

// AddAuthInfo adds a authentication
func (c *Config) AddAuthInfo(name string, auth *AuthInfo) {
	if c.AuthInfos == nil {
		c.AuthInfos = make(map[string]*AuthInfo)
	}
	c.AuthInfos[name] = auth
}

// HasValidProfile checks we have a current context
func (c *Config) HasValidProfile() error {
	if c.CurrentProfile == "" {
		return errors.New("no profile selected, please run korectl profile use <name>")
	}
	if !c.HasServer(c.GetCurrentProfile().Server) {
		return errors.New("profile does not have a server endpoint")
	}

	return nil
}

// HasProfile checks if the context exists in the config
func (c *Config) HasProfile(name string) bool {
	_, found := c.Profiles[name]

	return found
}

// HasServer checks if the context exists in the config
func (c *Config) HasServer(name string) bool {
	_, found := c.Servers[name]

	return found
}

// HasAuthInfo checks if the context exists in the config
func (c *Config) HasAuthInfo(name string) bool {
	_, found := c.AuthInfos[name]

	return found
}

// IsOIDCProviderConfigured checks if an OIDC provider is available
func (c *Config) IsOIDCProviderConfigured(name string) bool {
	info, found := c.AuthInfos[name]
	if !found {
		return false
	}

	return len(info.OIDC.ClientID) > 0 &&
		len(info.OIDC.ClientSecret) > 0 &&
		len(info.OIDC.AuthorizeURL) > 0
}

// HasSwagger checks if the swagger exists
func (c *Config) HasSwagger() (bool, error) {
	if _, err := os.Stat(c.GetSwaggerCachedFile()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// UpdateSwaggerCache updates the local swagger cache file
func (c *Config) UpdateSwaggerCache(content []byte) error {
	return ioutil.WriteFile(c.GetSwaggerCachedFile(), content, os.FileMode(0740))
}

// GetResourceChecksum requests the resource checksum
func (c *Config) GetResourceChecksum() (string, error) {
	v, err := c.request(c.GetResourceSumAPI())

	return string(v), err
}

// GetSwaggerChecksum requests the swagger checksum
func (c *Config) GetSwaggerChecksum() (string, error) {
	v, err := c.request(c.GetSwaggerSumAPI())

	return string(v), err
}

// GetSwagger returns the cached swagger
func (c *Config) GetSwagger() ([]byte, error) {
	return ioutil.ReadFile(c.GetSwaggerCachedFile())
}

// GetSwaggerFromAPI returns the swagger from api
func (c *Config) GetSwaggerFromAPI() ([]byte, error) {
	return c.request(c.GetSwaggerAPI())
}

// GetResourcesFromAPI returns the resource cache from api
func (c *Config) GetResourcesFromAPI() ([]byte, error) {
	return c.request(c.GetResourcesAPI())
}

func (c *Config) request(url string) ([]byte, error) {
	resp, err := hp.R().
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid response from apiserver: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// GetResourceSumAPI returns the api cache url
func (c *Config) GetResourceSumAPI() string {
	return fmt.Sprintf("%s/%s", c.GetAPI(), "classes/checksum")
}

// GetSwaggerSumAPI returns the location of the cached swagger checksum
func (c *Config) GetSwaggerSumAPI() string {
	return fmt.Sprintf("%s/%s", c.GetCurrentServer().Endpoint, "swagger.json?checksum=sha256")
}

// GetSwaggerAPI returns the api cache url
func (c *Config) GetSwaggerAPI() string {
	return fmt.Sprintf("%s/swagger.json", c.GetCurrentServer().Endpoint)
}

// GetResourcesAPI returns the api cache url
func (c *Config) GetResourcesAPI() string {
	return fmt.Sprintf("%s/classes", c.GetAPI())
}

// GetAPI returns the api server url
func (c *Config) GetAPI() string {
	return fmt.Sprintf("%s%s", c.GetCurrentServer().Endpoint, "/api/v1alpha1")
}

// RemoveServer removes a server instance
func (c *Config) RemoveServer(name string) {
	delete(c.Servers, name)
}

// RemoteUserInfo removes the user info
func (c *Config) RemoteUserInfo(name string) {
	delete(c.AuthInfos, name)
}

// RemoveProfile removes the profile
func (c *Config) RemoveProfile(name string) {
	p, found := c.Profiles[name]
	if !found {
		return
	}
	c.RemoveServer(p.Server)
	c.RemoteUserInfo(p.AuthInfo)

	delete(c.Profiles, name)
}

// Update writes the config to the korectl file
func (c *Config) Update() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	configPath := os.ExpandEnv(HubConfig)

	if err := os.MkdirAll(filepath.Dir(configPath), os.FileMode(0750)); err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, os.FileMode(0640))
}
