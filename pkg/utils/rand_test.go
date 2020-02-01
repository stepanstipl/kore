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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRand(t *testing.T) {
	u, err := Rand(12)
	require.NoError(t, err)
	assert.NotEmpty(t, u)
}

func TestRandLength(t *testing.T) {
	u, err := Rand(12)
	require.NoError(t, err)
	assert.Equal(t, 12, len(u))
}

func TestRandom(t *testing.T) {
	r := Random(12)
	assert.NotEmpty(t, r)
	assert.Equal(t, 12, len(r))
}

func TestRandomWithCharset(t *testing.T) {
	r := RandomWithCharset(12, DefaultiCharSet)
	assert.NotEmpty(t, r)
	assert.Equal(t, 12, len(r))
}

func TestRandomWithCharsetSet(t *testing.T) {
	r := RandomWithCharset(4, "a")
	assert.NotEmpty(t, r)
	assert.Equal(t, 4, len(r))
	assert.Equal(t, "aaaa", r)
}
