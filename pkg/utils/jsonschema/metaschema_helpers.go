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
	"encoding/json"
	"fmt"
)

// GetObjectIdentifier returns with the object field name which is marked as an identifier
// If the schema is not an object type or there is no identifier, it returns an empty string
func (m MetaSchemaDraft7Ext) GetObjectIdentifier() string {
	if m.Type == "object" {
		for name, prop := range m.Properties {
			if prop.Identifier {
				return name
			}
		}
	}
	return ""
}

func (m MetaSchemaDraft7Ext) GetSchemaForArrayItems(schema MetaSchemaDraft7Ext) MetaSchemaDraft7Ext {
	if schema.Type != "array" {
		panic("getSchemaForArrayItems should be called only for array schema definitions")
	}

	if schema.Items == nil {
		panic("array items schema is not set")
	}

	switch s := schema.Items.(type) {
	case map[string]interface{}:
		var res MetaSchemaDraft7Ext
		ref, ok := s["$ref"].(string)
		if ok {
			if ref == "#" {
				res = m
			} else {
				res = m.Definitions[ref]
			}
		} else {
			tmp, _ := json.Marshal(s)
			_ = json.Unmarshal(tmp, &res)
		}
		return res
	default:
		panic(fmt.Errorf("unexpected array items type: %T", schema.Items))
	}
}
