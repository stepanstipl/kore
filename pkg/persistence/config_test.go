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

package persistence

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
