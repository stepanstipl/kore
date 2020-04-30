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
)

func TestContainsOK(t *testing.T) {
	list := []string{"a", "b", "c"}

	assert.True(t, Contains("a", list))
	assert.True(t, Contains("c", list))
	assert.True(t, Contains("b", list))
}

func TestContainsBad(t *testing.T) {
	list := []string{"a", "b", "c"}

	assert.False(t, Contains("d", list))
}

func TestStringsSorted(t *testing.T) {
	a := []string{"a", "c", "b"}
	expected := []string{"a", "b", "c"}
	v := StringsSorted(a)
	assert.Equal(t, []string{"a", "c", "b"}, a)
	assert.Equal(t, expected, v)
}
