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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitUntilCompleteOK(t *testing.T) {
	timeout := 1 * time.Second
	interval := 1 * time.Second

	err := WaitUntilComplete(context.Background(), timeout, interval, func() (bool, error) {
		return true, nil
	})
	assert.NoError(t, err)
}

func TestWaitUntilCompleteErrored(t *testing.T) {
	timeout := 1 * time.Second
	interval := 1 * time.Second

	err := WaitUntilComplete(context.Background(), timeout, interval, func() (bool, error) {
		return false, errors.New("bad")
	})
	assert.Error(t, err)
}

func TestWaitUntilCompleteTimeout(t *testing.T) {
	timeout := 100 * time.Millisecond
	interval := 50 * time.Millisecond

	err := WaitUntilComplete(context.Background(), timeout, interval, func() (bool, error) {
		return false, nil
	})
	require.Error(t, err)
	assert.Equal(t, ErrCancelled, err)
}
