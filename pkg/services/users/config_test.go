/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigIsValidOK(t *testing.T) {
	config := Config{
		Driver:   "mysql",
		StoreURL: "noen",
	}
	assert.NoError(t, config.IsValid())
}

func TestConfigIsValidBad(t *testing.T) {
	cases := []struct {
		Config   Config
		Expected string
	}{
		{
			Config:   Config{},
			Expected: "no database driver configured",
		},
		{
			Config:   Config{Driver: "none", StoreURL: "dsd"},
			Expected: "unknown driver configured",
		},
		{
			Config:   Config{Driver: "mysql"},
			Expected: "no database url configured",
		},
	}
	for _, c := range cases {
		err := c.Config.IsValid()
		require.Error(t, err)
		assert.Equal(t, c.Expected, err.Error())
	}
}
