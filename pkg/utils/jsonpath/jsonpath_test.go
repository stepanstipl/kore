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

package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	document = `
{
	"apiVersion": "v1",
	"kind": "Test",
	"metadata": {
		"name": "test",
		"namespace": "test",
		"labels": {
			"one": "two"
		}
	},
	"items": ["one", "two", "three"],
	"spec": {
		"somthing": "yep"
	}
}
`
)

func TestIsValid(t *testing.T) {
	assert.True(t, IsValid(document))
}

func TestGet(t *testing.T) {
	v := Get(document, "metadata.namespace")
	require.NotNil(t, v)
	require.True(t, v.Exists())
	require.False(t, v.IsObject())
	require.Equal(t, "test", v.Value())
}

func TestGetJoin(t *testing.T) {
	v := Get(document, "items|@sjoin")
	require.NotNil(t, v)
	require.True(t, v.Exists())
	require.False(t, v.IsArray())
	require.Equal(t, "one,two,three", v.Value())
}
