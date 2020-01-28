/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearer(t *testing.T) {
	cases := []struct {
		Header   string
		Expected string
		Found    bool
	}{
		{},
		{Header: "nothing"},
		{Header: "Bearer"},
		{Header: "Bearer test", Expected: "test", Found: true},
		{Header: "bearer test", Expected: "test", Found: true},
	}
	for _, c := range cases {
		v, found := GetBearerToken(c.Header)
		assert.Equal(t, c.Found, found)
		assert.Equal(t, c.Expected, v)
	}
}
