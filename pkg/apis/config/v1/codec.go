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

package v1

import (
	"encoding/base64"
	"fmt"
)

// Encode is responsible for ensuring the secret is encoded
func (s *Secret) Encode() *Secret {
	for k, v := range s.Spec.Data {
		s.Spec.Data[k] = base64.RawStdEncoding.EncodeToString([]byte(v))
	}

	return s
}

// Decode is responsible for decoding the secret
func (s *Secret) Decode() error {
	for k, v := range s.Spec.Data {
		decoded, err := base64.RawStdEncoding.DecodeString(v)
		if err != nil {
			return fmt.Errorf("key: %s, value: %s is not base64 encoded", k, v)
		}
		s.Spec.Data[k] = string(decoded)
	}

	return nil
}
