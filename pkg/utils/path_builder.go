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
	"fmt"
	"strings"
)

// PathBuilder provides a helper method for builder url paths
type PathBuilder struct {
	base string
}

// NewPathBuilder returns a new builder
func NewPathBuilder(base string) PathBuilder {
	trimmed := strings.TrimSuffix(base, "/")

	return PathBuilder{base: trimmed}
}

// Add a path buidler on top of
func (b PathBuilder) Add(v string) PathBuilder {
	return NewPathBuilder(b.Path(v))
}

// Base returns the base path
func (b PathBuilder) Base() string {
	return b.base
}

// Path returns the base plus the path
func (b PathBuilder) Path(p string) string {
	v := fmt.Sprintf("%s/%s", b.base, strings.TrimPrefix(p, "/"))

	return strings.TrimSuffix(v, "/")
}
