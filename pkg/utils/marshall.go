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

package utils

import (
	"bytes"
	"encoding/json"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/runtime"
)

// EncodeRuntimeObjectToYAML is used to encode the object to a yaml document
func EncodeRuntimeObjectToYAML(object runtime.Object) ([]byte, error) {
	b := &bytes.Buffer{}

	// @step: encode to json first of all
	if err := json.NewEncoder(b).Encode(object); err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(b.Bytes())
}
