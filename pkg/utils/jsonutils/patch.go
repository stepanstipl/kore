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

package jsonutils

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
)

func Diff(v1, v2 interface{}) ([]byte, error) {
	j1, err := json.Marshal(v1)
	if err != nil {
		return nil, err
	}

	j2, err := json.Marshal(v2)
	if err != nil {
		return nil, err
	}

	return jsonpatch.CreateMergePatch(j1, j2)
}
