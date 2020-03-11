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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDurationPtr(t *testing.T) {
	d := DurationPtr(1 * time.Minute)
	require.NotNil(t, d)
	assert.Equal(t, time.Minute*1, *d)
}

func TestStringPtr(t *testing.T) {
	s := StringPtr("hello")
	require.NotNil(t, s)
	assert.Equal(t, "hello", *s)
}
