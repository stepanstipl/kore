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

func TestRetryContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := Retry(ctx, 0, true, 1000*time.Millisecond, func() (bool, error) {
		return false, nil
	})
	assert.Equal(t, ErrCancelled, err)
}

func TestRetryWithSuccess(t *testing.T) {
	err := Retry(context.Background(), 0, true, 10*time.Millisecond, func() (bool, error) {
		return true, nil
	})
	assert.NoError(t, err)
}

func TestRetryMaxAttempts(t *testing.T) {
	var list []string

	err := Retry(context.Background(), 3, true, 10*time.Millisecond, func() (bool, error) {
		list = append(list, "done")

		return false, nil
	})
	require.Error(t, err)
	assert.Equal(t, ErrReachMaxAttempts, err)
}

func TestRetryWithError(t *testing.T) {
	err := Retry(context.Background(), 3, true, 10*time.Millisecond, func() (bool, error) {
		return false, errors.New("bad")
	})
	require.Error(t, err)
	assert.Equal(t, "bad", err.Error())
}

func TestRetryWithTimeout(t *testing.T) {
	err := RetryWithTimeout(context.Background(), 50*time.Millisecond, 30*time.Millisecond, func() (bool, error) {
		return false, nil
	})
	require.Error(t, err)
	assert.Equal(t, ErrCancelled, err)
}
