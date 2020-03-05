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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSleep(t *testing.T) {
	assert.False(t, Sleep(context.Background(), 1*time.Millisecond))
}

func TestSleepWithCancel(t *testing.T) {
	c, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	assert.True(t, Sleep(c, 10*time.Second))
}
