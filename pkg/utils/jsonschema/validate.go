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
