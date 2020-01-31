/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
