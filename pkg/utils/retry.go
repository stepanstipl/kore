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
	"time"

	"github.com/jpillora/backoff"
)

var (
	// ErrReachMaxAttempts indicates we hit the limit
	ErrReachMaxAttempts = errors.New("reached max attempts")
)

const (
	// MaxAttempts is the max attempts
	MaxAttempts = 99999999
)

// RetryFunc performs the operation
type RetryFunc func() (bool, error)

// RetryWithTimeout creates a retry with a specific timeout
func RetryWithTimeout(ctx context.Context, timeout, interval time.Duration, retryFn RetryFunc) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return Retry(ctx, 0, true, interval, retryFn)
}

// Retry is used to retry an operation multiple times under a context
func Retry(ctx context.Context, attempts int, jitter bool, minInterval time.Duration, retryFn RetryFunc) error {
	// @hack: quick way to do this for now
	if attempts == 0 {
		attempts = MaxAttempts
	}

	backoff := &backoff.Backoff{
		Min:    minInterval,
		Max:    minInterval * 2,
		Factor: 1,
		Jitter: jitter,
	}

	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return ErrCancelled
		default:
		}

		finished, err := retryFn()
		if err != nil {
			return err
		}
		if finished {
			return nil
		}

		Sleep(ctx, backoff.Duration())
	}

	return ErrReachMaxAttempts
}
