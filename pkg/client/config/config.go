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

package config

import (
	"errors"
	"io"
	"os"

	"github.com/appvia/kore/pkg/version"

	"gopkg.in/yaml.v2"
)

// New creates a configuration
func New(reader io.Reader) (*Config, error) {
	config := &Config{}

	return config, yaml.NewDecoder(reader).Decode(config)
}

// NewEmpty returns an empty configuration
func NewEmpty() *Config {
	return &Config{
		AuthInfos: make(map[string]*AuthInfo),
		Profiles:  make(map[string]*Profile),
		Servers:   make(map[string]*Server),
		Version:   version.Release,
	}
}

// NewFromFile returns the the configuration from a file
func NewFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return New(file)
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() error {
	return nil
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

// ListProfiles returns a list of profile names
func (c *Config) ListProfiles() []string {
	if c.Profiles == nil {
		return nil
	}
	var list []string

	for k := range c.Profiles {
		list = append(list, k)
	}

	return list
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

// AddProfile adds a profile to the config
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
		return errors.New("no profile selected")
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

// RemoveServer removes a server instance
func (c *Config) RemoveServer(name string) {
	delete(c.Servers, name)
}

// RemoveUserInfo removes the user info
func (c *Config) RemoveUserInfo(name string) {
	delete(c.AuthInfos, name)
}

// RemoveProfile removes the profile
func (c *Config) RemoveProfile(name string) {
	p, found := c.Profiles[name]
	if !found {
		return
	}
	c.RemoveServer(p.Server)
	c.RemoveUserInfo(p.AuthInfo)

	delete(c.Profiles, name)
}

// Update writes the config to the file
func (c *Config) Update(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(c)
}
