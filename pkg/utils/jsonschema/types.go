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

package jsonschema

import (
	"fmt"
	"strconv"
	"strings"
)

// Schema is the JSON schema object
type Schema struct {
	Properties map[string]*Property `json:"properties"`
}

// Property is an object property
type Property struct {
	Default     interface{}   `json:"default"`
	Const       interface{}   `json:"const"`
	Description string        `json:"description"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	Enum        []interface{} `json:"enum"`
}

func (p Property) ParseConst() (interface{}, error) {
	return p.parseValue(p.Const)
}

func (p Property) ParseDefault() (interface{}, error) {
	return p.parseValue(p.Default)
}
func (p Property) parseValue(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	// Check if the value is a complete template expression ("{{ ... }}") and return it without type checking
	if val, ok := value.(string); ok {
		if strings.HasPrefix(val, "{{") && strings.HasSuffix(val, "}}") {
			return val, nil
		}
	}

	switch p.Type {
	case "string":
		if val, ok := value.(string); ok {
			return val, nil
		}
		return fmt.Sprintf("%v", value), nil
	case "boolean":
		if val, ok := value.(bool); ok {
			return val, nil
		}
		val, err := strconv.ParseBool(fmt.Sprintf("%v", value))
		return val, err
	case "integer":
		if val, ok := value.(int64); ok {
			return val, nil
		}
		if val, ok := value.(float64); ok {
			return int64(val), nil
		}
		val, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
		return val, err
	case "number":
		if val, ok := value.(int64); ok {
			return val, nil
		}
		if val, ok := value.(float64); ok {
			return val, nil
		}
		val, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		return val, err
	default:
		return value, nil
	}
}
