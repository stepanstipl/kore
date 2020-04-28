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
	"strings"

	"github.com/appvia/kore/pkg/utils/validation"

	"github.com/xeipuuv/gojsonschema"
)

// Validate runs a JSON schema validation using the given schema against the passed object
func Validate(schemaJSON string, subject string, data interface{}) error {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaJSON))
	if err != nil {
		panic(fmt.Errorf("failed to compile plan schema: %v", err))
	}

	var loader gojsonschema.JSONLoader
	switch d := data.(type) {
	case []byte:
		loader = gojsonschema.NewBytesLoader(d)
	case string:
		loader = gojsonschema.NewStringLoader(d)
	default:
		loader = gojsonschema.NewGoLoader(d)
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
