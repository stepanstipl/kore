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

package application

import (
	"strconv"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"sigs.k8s.io/yaml"
)

type YAMLMap map[string]interface{}

func (y YAMLMap) MarshalJSON() ([]byte, error) {
	if len(y) == 0 {
		return []byte(`""`), nil
	}

	yamlData, err := yaml.Marshal(map[string]interface{}(y))
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Quote(string(yamlData))), nil
}

func (y *YAMLMap) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	raw, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if raw == "" {
		return nil
	}

	res := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(raw), &res); err != nil {
		return err
	}

	*y = res

	return nil
}

type ProviderData struct {
	Resources []corev1.Ownership `json:"resources,omitempty"`
}
