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
