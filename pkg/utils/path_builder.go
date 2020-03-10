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
