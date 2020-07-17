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
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConfig = `
current-profile: local
profiles:
  local:
    server: local
    team: kore
    user: local
servers:
  local:
    server: "http://127.0.0.1:10080"
users:
  local:
    oidc:
      access-token: test
      authorize-url: test
      client-id: test
      client-secret: test
      id-token: test
      refresh-token: test
      token-url: test
`
)

func makeTestConfig() io.Reader {
	return strings.NewReader(testConfig)
}

func newTestConfig(t *testing.T) *Config {
	c, err := New(makeTestConfig())
	require.NoError(t, err)
	require.NotNil(t, c)

	return c
}

func TestNewEmpty(t *testing.T) {
	c := NewEmpty()
	assert.NotNil(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(makeTestConfig())
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestCurrentProfile(t *testing.T) {
	c := newTestConfig(t)
	assert.Equal(t, "local", c.CurrentProfile)
}

func TestGetProfile(t *testing.T) {
	c := newTestConfig(t)
	profile := c.GetProfile("local")

	require.NotNil(t, profile)
	assert.Equal(t, "local", profile.Server)
	assert.Equal(t, "local", profile.AuthInfo)
	assert.Equal(t, "kore", profile.Team)
}

func TestGetCurrentServer(t *testing.T) {
	c := newTestConfig(t)
	current := c.GetServer("local")

	require.NotNil(t, current)
	assert.Equal(t, "http://127.0.0.1:10080", current.Endpoint)
}

func TestGetCurrentAuthInfo(t *testing.T) {
	c := newTestConfig(t)
	current := c.GetAuthInfo("local")

	require.NotNil(t, current)
	require.NotNil(t, current.OIDC)
	assert.Equal(t, "test", current.OIDC.AccessToken)
	assert.Equal(t, "test", current.OIDC.AuthorizeURL)
	assert.Equal(t, "test", current.OIDC.ClientID)
	assert.Equal(t, "test", current.OIDC.ClientSecret)
	assert.Equal(t, "test", current.OIDC.IDToken)
	assert.Equal(t, "test", current.OIDC.RefreshToken)
	assert.Equal(t, "test", current.OIDC.TokenURL)
}

func TestAddProfile(t *testing.T) {
	c := newTestConfig(t)
	c.AddProfile("demo", &Profile{
		AuthInfo: "demo",
		Team:     "demo",
		Server:   "demo",
	})
	require.True(t, c.HasProfile("demo"))
}

func TestIsValid(t *testing.T) {
	c := newTestConfig(t)
	assert.Nil(t, c.IsValid())
}
