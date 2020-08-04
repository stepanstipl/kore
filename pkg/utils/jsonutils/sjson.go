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

package jsonutils

import (
	"strconv"
	"strings"

	"github.com/tidwall/sjson"
)

// SetJSONProperty sets a json property
func SetJSONProperty(document []byte, key, value string) ([]byte, error) {
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return sjson.SetRawBytes(document, key, []byte(value))
	}
	return sjson.SetBytes(document, key, func(v string) interface{} {
		if val, err := strconv.ParseBool(v); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
		if val, err := strconv.Unquote(v); err == nil {
			return val
		}

		return v
	}(value))
}
