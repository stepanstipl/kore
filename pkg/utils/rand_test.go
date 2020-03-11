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
