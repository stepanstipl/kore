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
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// Validate runs a JSON schema validation using the given schema against the passed object
func Validate(schemaJSON string, subject string, data interface{}) error {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaJSON))
	if err != nil {
		panic(fmt.Errorf("failed to compile gkePlanSchema: %v", err))
	}

	res, err := schema.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return fmt.Errorf("failed to parse data for validation: %s", err)
	}
	if !res.Valid() {
		errStr := fmt.Sprintf("%s has failed validation:\n", subject)
		for _, err := range res.Errors() {
			errStr += fmt.Sprintf(" * %s: %s\n", err.Field(), err.Description())
		}
		return errors.New(errStr)
	}

	return nil
}
