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

package utils

import (
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClaims(t *testing.T) {
	claims := jwt.MapClaims{}

	c := NewClaims(claims)
	assert.NotEmpty(t, c)
}

func TestClaimsGetStringOK(t *testing.T) {
	claims := jwt.MapClaims{
		"hello": "world",
	}
	c := NewClaims(claims)
	require.NotEmpty(t, c)

	v, found := c.GetString("not_there")
	assert.Empty(t, v)
	assert.False(t, found)

	v, found = c.GetString("hello")
	assert.Equal(t, "world", v)
	assert.True(t, found)
}

func TestClaimsGetStringNotType(t *testing.T) {
	claims := jwt.MapClaims{
		"hello": 1,
	}
	c := NewClaims(claims)
	require.NotEmpty(t, c)

	v, found := c.GetString("hello")
	assert.Empty(t, v)
	assert.False(t, found)
}
