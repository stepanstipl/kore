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

package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	s, err := FromString("hello")()
	require.NoError(t, err)
	assert.Equal(t, "hello", s)
}

func TestFromBytes(t *testing.T) {
	s, err := FromBytes([]byte("hello"))()
	require.NoError(t, err)
	assert.Equal(t, "hello", s)
}

func TestFromReaderOK(t *testing.T) {
	s, err := FromReader(strings.NewReader("hello"))()
	require.NoError(t, err)
	assert.Equal(t, "hello", s)
}

func TestFromReaderBad(t *testing.T) {
	s, err := FromReader(nil)()
	require.Error(t, err)
	assert.Equal(t, "", s)
	assert.Equal(t, err, ErrInvalidReader)
}

func TestColumnOK(t *testing.T) {
	c, err := Column("test", "test")()
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestColumnNoName(t *testing.T) {
	c, err := Column("", "test")()
	require.Error(t, err)
	assert.Nil(t, c)
}

func TestColumnNilMethod(t *testing.T) {
	c, err := Column("test", "test", nil)()
	require.Error(t, err)
	assert.Nil(t, c)
}
