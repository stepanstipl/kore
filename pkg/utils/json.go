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

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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

func ApiExtJSONEquals(j1, j2 *apiextv1.JSON) bool {
	if ApiExtJSONEmpty(j1) && ApiExtJSONEmpty(j2) {
		return true
	}

	if ApiExtJSONEmpty(j1) || ApiExtJSONEmpty(j2) {
		return false
	}

	return bytes.Equal(j1.Raw, j2.Raw)
}

func ApiExtJSONEmpty(j *apiextv1.JSON) bool {
	if j == nil {
		return true
	}

	return len(j.Raw) == 0 || bytes.Equal(j.Raw, []byte("{}"))
}
