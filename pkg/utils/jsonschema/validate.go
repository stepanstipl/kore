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
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/appvia/kore/pkg/utils/validation"

	"github.com/xeipuuv/gojsonschema"
)

var emptyObjectLoader = gojsonschema.NewStringLoader("{}")

// Validate runs a JSON schema validation using the given schema against the passed object
func Validate(schemaJSON string, subject string, data interface{}) error {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaJSON))
	if err != nil {
		return fmt.Errorf("failed to compile schema: %v", err)
	}

	loader := emptyObjectLoader

	switch d := data.(type) {
	case apiextv1.JSON:
		if len(d.Raw) > 0 {
			loader = gojsonschema.NewBytesLoader(d.Raw)
		}
	case *apiextv1.JSON:
		if d != nil && len(d.Raw) > 0 {
			loader = gojsonschema.NewBytesLoader(d.Raw)
		}
	case []byte:
		if len(d) > 0 {
			loader = gojsonschema.NewBytesLoader(d)
		}
	case string:
		if len(d) > 0 {
			loader = gojsonschema.NewStringLoader(d)
		}
	default:
		if d != nil {
			loader = gojsonschema.NewGoLoader(d)
		}
	}

	res, err := schema.Validate(loader)
	if err != nil {
		return fmt.Errorf("%s has failed validation: %w", subject, err)
	}
	if !res.Valid() {
		ve := validation.NewError("%s has failed validation", subject)
		for _, err := range res.Errors() {
			switch err.(type) {
			case *gojsonschema.ConditionElseError, *gojsonschema.ConditionThenError:
				// Ignore these errors
			default:
				field := err.Field()
				// in the case of required error type, get the field name from the message "fieldName is required"
				// use this for "field" instead of "(root)"
				if err.Type() == "required" && field == validation.FieldRoot {
					field = strings.Fields(err.Description())[0]
				}
				ve.AddFieldError(field, validation.ErrorCode(err.Type()), err.Description())
			}
		}
		return ve
	}

	return nil
}

// ValidateImmutableProperties validates recursively all immutable fields
// The `fieldPrefix` parameter will be added as a prefix to the validation error field names.
//
// If your schema contains an array of objects, and the objects have a unique identifier field (e.g. `name`),
// you should set the `identifier` flag on the identifier field, so the validation can always compare the same objects.
func ValidateImmutableProperties(schemaJSON string, subject, fieldPrefix string, j1, j2 []byte) error {
	v, err := newValidator(schemaJSON, subject, fieldPrefix)
	if err != nil {
		return err
	}

	return v.validateImmutableProperties(
		gjson.GetBytes(j1, "@this"),
		gjson.GetBytes(j2, "@this"),
	)
}

type validator struct {
	schema      MetaSchemaDraft7Ext
	fieldPrefix string
	ve          *validation.Error
}

func newValidator(schemaBytes string, subject, fieldPrefix string) (*validator, error) {
	schema := MetaSchemaDraft7Ext{}
	if err := json.Unmarshal([]byte(schemaBytes), &schema); err != nil {
		return nil, fmt.Errorf("schema is invalid: %w", err)
	}

	return &validator{
		schema:      schema,
		fieldPrefix: fieldPrefix,
		ve:          validation.NewError("%s has failed validation", subject),
	}, nil
}

func (v *validator) validateImmutableProperties(v1, v2 gjson.Result) error {
	v.validate(v.schema, "", v1, v2)
	if v.ve.HasErrors() {
		return v.ve
	}

	return nil
}

func (v *validator) validate(schema MetaSchemaDraft7Ext, path string, v1, v2 gjson.Result) {
	if !v2.Exists() {
		return
	}

	if schema.Immutable && !reflect.DeepEqual(v.normalize(v1.Value()), v.normalize(v2.Value())) {
		v.ve.AddFieldError(pathJoin(v.fieldPrefix, path), validation.ReadOnly, "updating the field is not allowed")
		return
	}

	switch schema.Type {
	case "object":
		for name, prop := range schema.Properties {
			v.validate(prop, pathJoin(path, name), v1.Get(name), v2.Get(name))
		}
	case "array":
		j1a, ok := v1.Value().([]interface{})
		if ok {
			itemsSchema := v.schema.GetSchemaForArrayItems(schema)
			idField := itemsSchema.GetObjectIdentifier()
			if idField != "" {
				for i := range j1a {
					id := v1.Get(fmt.Sprintf("%d.%s", i, idField))
					v.validate(
						itemsSchema,
						pathJoin(path, strconv.Itoa(i)),
						v1.Get(strconv.Itoa(i)),
						v2.Get(fmt.Sprintf("#(%s==%q)", idField, id)),
					)
				}
			} else {
				for i := range j1a {
					v.validate(
						itemsSchema,
						pathJoin(path, strconv.Itoa(i)),
						v1.Get(strconv.Itoa(i)),
						v2.Get(strconv.Itoa(i)),
					)
				}
			}
		}
	}
}

func (v *validator) normalize(value interface{}) interface{} {
	switch vt := value.(type) {
	case []interface{}:
		if len(vt) == 0 {
			return nil
		}
		for i, e := range vt {
			vt[i] = v.normalize(e)
		}
		return vt
	case map[string]interface{}:
		if len(vt) == 0 {
			return nil
		}
		for k, e := range vt {
			vt[k] = v.normalize(e)
		}
		return vt
	default:
		return value
	}
}

func pathJoin(v ...string) string {
	tmp := make([]string, 0, len(v))
	for _, e := range v {
		if e != "" {
			tmp = append(tmp, e)
		}
	}
	return strings.Join(tmp, ".")
}
