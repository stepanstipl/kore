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

package render

import (
	"strings"
	"time"
)

// Age formats to an age
func Age() PrinterColumnFormatter {
	return func(value string) string {
		created, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return "Invalid"
		}

		return strings.Split(time.Since(created).String(), ".")[0] + "s"
	}
}

// Default used the v if no value set
func Default(v string) PrinterColumnFormatter {
	return func(value string) string {
		if value == "" {
			return v
		}

		return value
	}
}

// IfEqual checks if v is equal and if so returns x
func IfEqual(v, x string) PrinterColumnFormatter {
	return func(value string) string {
		if value == x {
			return x
		}

		return value
	}
}

// IfEqualOr checks if v is equal and if so returns x
func IfEqualOr(v, a, b string) PrinterColumnFormatter {
	return func(value string) string {
		if value == v {
			return a
		}

		return b
	}
}
