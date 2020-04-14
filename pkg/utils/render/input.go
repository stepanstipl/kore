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

package render

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Column is just syntaxtic sugar
func Column(name, path string, formatters ...PrinterColumnFormatter) PrinterColumnFunc {
	return func() (*column, error) {
		if name == "" {
			return nil, &ErrInvalidColumn{"no name"}
		}
		for _, x := range formatters {
			if x == nil {
				return nil, &ErrInvalidColumn{"formatter method is nil"}
			}
		}

		return &column{name: name, path: path, formatters: formatters}, nil
	}
}

// FromStruct reads the resource from a struct
func FromStruct(v interface{}) ResourceInputFunc {
	return func() (string, error) {
		b := &bytes.Buffer{}
		if err := json.NewEncoder(b).Encode(v); err != nil {
			return "", err
		}

		return b.String(), nil
	}
}

// FromBytes reads the resource from a bytes slice
func FromBytes(v []byte) ResourceInputFunc {
	return func() (string, error) {
		return string(v), nil
	}
}

// FromString reads from a string
func FromString(v string) ResourceInputFunc {
	return func() (string, error) {
		return v, nil
	}
}

// FromReader reads the resource from a io.reader
func FromReader(v io.Reader) ResourceInputFunc {
	return func() (string, error) {
		if v == nil {
			return "", ErrInvalidReader
		}

		data, err := ioutil.ReadAll(v)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}
}
