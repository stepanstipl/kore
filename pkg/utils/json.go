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
	"io"
)

// EncodeToJSON encodes the struct to json
func EncodeToJSON(in interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(in); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

// DecodeToJSON encodes the struct to json
func DecodeToJSON(data io.Reader, in interface{}) error {
	if err := json.NewDecoder(data).Decode(in); err != nil {
		return err
	}

	return nil
}
