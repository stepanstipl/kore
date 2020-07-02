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
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Compact compacts the JSON data
// Object keys will be sorted alphabetically
func Compact(j []byte) ([]byte, error) {
	j = bytes.TrimSpace(j)

	if !bytes.HasPrefix(j, []byte("{")) {
		buf := bytes.NewBuffer(make([]byte, 0, len(j)))
		if err := json.Compact(buf, j); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(j, &data); err != nil {
		return nil, err
	}

	return json.Marshal(wrapMap(data))
}

func wrapMap(m map[string]interface{}) sortedMap {
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			m[k] = wrapMap(vm)
		}
	}
	return sortedMap(m)
}

type sortedMap map[string]interface{}

func (s sortedMap) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(s))

	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := bytes.NewBuffer(make([]byte, 0, 128))
	buf.WriteRune('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteRune(',')
		}
		_, _ = fmt.Fprintf(buf, "%q:", k)
		val, err := json.Marshal(s[k])
		if err != nil {
			return nil, err
		}
		_, _ = buf.Write(val)
	}
	buf.WriteRune('}')

	return buf.Bytes(), nil
}
