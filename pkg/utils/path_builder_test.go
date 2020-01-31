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

func TestNewPathBuilder(t *testing.T) {
	b := NewPathBuilder("/base/")
	assert.NotNil(t, b)
}

func TestNewPathBuilderBase(t *testing.T) {
	b := NewPathBuilder("/base")
	require.NotNil(t, b)
	assert.Equal(t, "/base", b.Base())
}

func TestPathBuilderURL(t *testing.T) {
	b := NewPathBuilder("/base")
	require.NotNil(t, b)
	assert.Equal(t, "/base/v1/teams", b.Path("v1/teams"))
}

func TestPathBuilderAdd(t *testing.T) {
	b := NewPathBuilder("/base")
	require.NotNil(t, b)
	ba := b.Add("teams")

	assert.Equal(t, "/base/teams", ba.Base())
}
