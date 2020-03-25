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

package jsonschema_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/appvia/kore/pkg/utils/jsonschema"
)

const testSchema = `
{
	"description": "GKE Cluster Plan Schema",
	"type": "object",
	"additionalProperties": false,
	"required": [
		"p1",
		"p2"
	],
	"properties": {
		"p1": {
			"type": "number",
			"multipleOf": 1,
			"minimum": 10,
			"maximum": 20
		},
		"p2": {
			"type": "string",
			"minLength": 1
		},
		"p3": {
			"type": "boolean"
		},
		"p4": {
			"type": "array",
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": [
					"p41",
					"p42"
				],
				"properties": {
					"p41": {
						"type": "string"
					},
					"p42": {
						"type": "number"
					}
				}
			},
			"minItems": 1
		}
	}
}
`

const minimumOverride = `
{
	"minimum": 100,
	"maximum": 200
}
`

const recursiveOverride = `
{
	"items": {
		"properties": {
			"p41": {
				"minLength": 10
			},
			"p42": {
				"minimum": 20
			}
		}
	}
}
`

const invalidOverride = `
{
	"type": "foo"
}
`

func TestNewSchema(t *testing.T) {
	_, err := jsonschema.NewSchemaFromString(testSchema)
	require.NoError(t, err)
}

func TestPropertyOverride(t *testing.T) {
	schema, err := jsonschema.NewSchemaFromString(testSchema)
	require.NoError(t, err)

	err = schema.UpdateProperty("p1", minimumOverride)
	assert.NoError(t, err)

	properties := parseSchema(schema)["properties"].(map[string]interface{})
	p1 := properties["p1"].(map[string]interface{})

	require.Equal(t, 100.0, p1["minimum"])
	require.Equal(t, 200.0, p1["maximum"])
}

func TestRecursiveOverride(t *testing.T) {
	schema, err := jsonschema.NewSchemaFromString(testSchema)
	require.NoError(t, err)

	err = schema.UpdateProperty("p4", recursiveOverride)
	assert.NoError(t, err)

	properties := parseSchema(schema)["properties"].(map[string]interface{})
	p4 := properties["p4"].(map[string]interface{})
	p4Items := p4["items"].(map[string]interface{})
	p4Properties := p4Items["properties"].(map[string]interface{})
	p41 := p4Properties["p41"].(map[string]interface{})
	p42 := p4Properties["p42"].(map[string]interface{})

	assert.Equal(t, 10.0, p41["minLength"])
	assert.Equal(t, 20.0, p42["minimum"])
}

func TestSchemaMustStayValid(t *testing.T) {
	schema, err := jsonschema.NewSchemaFromString(testSchema)
	require.NoError(t, err)

	err = schema.UpdateProperty("p1", invalidOverride)
	assert.Error(t, err)
}

func parseSchema(schema *jsonschema.Schema) map[string]interface{} {
	var document map[string]interface{}
	_ = json.Unmarshal([]byte(schema.String()), &document)
	return document
}
