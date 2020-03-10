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

// ToPlural converts the type to a plural
func ToPlural(name string) string {
	if strings.HasSuffix(name, "ss") {
		return fmt.Sprintf("%ses", name)
	}
	if strings.HasSuffix(name, "ys") {
		return fmt.Sprintf("%sies", strings.TrimSuffix(name, "ys"))
	}
	if strings.HasSuffix(name, "es") || strings.HasSuffix(name, "s") {
		return name
	}

	return fmt.Sprintf("%ss", name)
}
