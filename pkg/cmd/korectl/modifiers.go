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

package korectl

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// RegisterModifiers is called to register the jsonpath modifiers
func RegisterModifiers() {
	gjson.AddModifier("sjoin", func(value, arg string) string {
		delimiter := arg
		if delimiter == "" {
			delimiter = ","
		}

		// @step: ensure this is an array - else return the value
		if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
			return value
		}

		// @step: decode the array into a slice
		var items []interface{}
		if err := json.NewDecoder(strings.NewReader(value)).Decode(&items); err != nil {
			return value
		}

		// @step: append and join the values together
		var list []string
		for _, x := range items {
			list = append(list, fmt.Sprintf("%v", x))
		}

		return fmt.Sprintf("\"%s\"", strings.Join(list, delimiter))
	})
}
