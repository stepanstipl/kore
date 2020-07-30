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
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// YAMLDocumentsFromString returns a series of documents from the string
func YAMLDocumentsFromString(content string) ([]string, error) {
	return YAMLDocuments(strings.NewReader(content))
}

// YAMLDocuments returns a collection of documents from the reader
func YAMLDocuments(reader io.Reader) ([]string, error) {
	// @step: read in the content of the file
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	splitter := regexp.MustCompile("(?m)^---\n")

	return splitter.Split(string(content), -1), nil
}

// ToYAML marshalls the struct to yaml
func ToYAML(in interface{}) ([]byte, error) {
	b := &bytes.Buffer{}

	err := yaml.NewEncoder(b).Encode(in)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
