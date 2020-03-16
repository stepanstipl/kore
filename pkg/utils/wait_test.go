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
	assert.Equal(t, ErrTimeout, err)
}
