/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("operation has timed out")
)

// WaitUntilComplete calls the condition on every interval and check for true, nil. An error indicates a hard error and exits
func WaitUntilComplete(ctx context.Context, timeout time.Duration, interval time.Duration, condition func() (bool, error)) error {
	nc, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if interval <= 0 {
		return errors.New("invalid interval")
	}

	for {
		select {
		case <-nc.Done():
			return ErrTimeout
		default:
		}

		if done, err := condition(); err != nil {
			return err
		} else if done {
			return nil
		}

		time.Sleep(interval)
	}
}
