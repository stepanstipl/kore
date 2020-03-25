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
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type Schema struct {
	document map[string]interface{}
	schema   *gojsonschema.Schema
}

func NewSchemaFromString(input string) (*Schema, error) {
	var document map[string]interface{}
	if err := json.Unmarshal([]byte(input), &document); err != nil {
		return nil, err
	}
	_, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(document))
	if err != nil {
		return nil, fmt.Errorf("failed to compile JSON schema: %s", err)
	}

	s := &Schema{
		document: document,
	}

	if err := s.parse(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Schema) UpdateProperty(name string, propertySchema string) error {
	var s2 map[string]interface{}
	if err := json.Unmarshal([]byte(propertySchema), &s2); err != nil {
		return err
	}
	properties := s.document["properties"].(map[string]interface{})
	s1, ok := properties[name]
	if !ok {
		return fmt.Errorf("%q property does not exist", name)
	}

	if err := mergeSchema(s1.(map[string]interface{}), s2); err != nil {
		return err
	}

	if err := s.parse(); err != nil {
		return err
	}

	return nil
}

func (s *Schema) Validate(subject string, data interface{}) error {
	res, err := s.schema.Validate(gojsonschema.NewGoLoader(data))
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

func (s *Schema) parse() error {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(s.document))
	if err != nil {
		return err
	}
	s.schema = schema
	return nil
}

func (s *Schema) String() string {
	res, _ := json.Marshal(s.document)
	return string(res)
}

func mergeSchema(s1, s2 map[string]interface{}) error {
	for k2, v2 := range s2 {
		switch k2 {
		case "items":
			if v1, ok := s1[k2]; ok {
				if err := mergeSchema(v1.(map[string]interface{}), v2.(map[string]interface{})); err != nil {
					return err
				}
			} else {
				return errors.New("\"items\" attribute can not be set on a non-array type")
			}
		case "properties":
			if v1, ok := s1[k2]; ok {
				if err := mergeProperties(v1.(map[string]interface{}), v2.(map[string]interface{})); err != nil {
					return err
				}
			} else {
				return errors.New("\"properties\" attribute can not be set on a non-object type")
			}
		default:
			s1[k2] = v2
		}
	}
	return nil
}

func mergeProperties(s1, s2 map[string]interface{}) error {
	for k2, v2 := range s2 {
		v1, ok := s1[k2]
		if !ok {
			return fmt.Errorf("%q property does not exist", k2)
		}
		if err := mergeSchema(v1.(map[string]interface{}), v2.(map[string]interface{})); err != nil {
			return err
		}
	}
	return nil
}
