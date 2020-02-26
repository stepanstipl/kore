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

// GetCurrentServer returns the server in the context
func (c *Config) GetCurrentServer() *Server {
	ct := c.Contexts[c.CurrentContext]
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
	ct := c.Contexts[c.CurrentContext]
	if ct == nil {
		return &AuthInfo{}
	}

	a := c.AuthInfos[ct.AuthInfo]

	if a == nil {
		return &AuthInfo{}
	}

	return a
}

// HasContext checks if the context exists in the config
func (c *Config) HasContext(name string) bool {
	_, found := c.Contexts[name]

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

func (c *Config) UpdateSwaggerCache(content []byte) error {
	return ioutil.WriteFile(c.GetSwaggerCachedFile(), content, os.FileMode(0740))
}

func (c *Config) GetResourceChecksum() (string, error) {
	v, err := c.request(c.GetResourceSumAPI())

	return string(v), err
}

func (c *Config) GetSwaggerChecksum() (string, error) {
	v, err := c.request(c.GetSwaggerSumAPI())

	return string(v), err
}

func (c *Config) GetSwagger() ([]byte, error) {
	return ioutil.ReadFile(c.GetSwaggerCachedFile())
}

func (c *Config) GetSwaggerFromAPI() ([]byte, error) {
	return c.request(c.GetSwaggerAPI())
}

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

// GetResourceCacheAPI returns the api cache url
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
