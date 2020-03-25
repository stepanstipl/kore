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
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadFileOrStdin is responsible for reading from the file or stdin
func ReadFileOrStdin(path string) ([]byte, error) {
	if path == "-" {
		return ioutil.ReadAll(os.Stdin)
	}

	return ioutil.ReadFile(path)
}

// FileExists checks if a file exists
func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return !info.IsDir(), nil
}

// EnsureFileExists creates an (empty) file if filename does not exist (returning true)
// or does nothing if it already exists (returning false). Will create missing directories.
func EnsureFileExists(filename string) (bool, error) {
	exists, err := FileExists(filename)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return false, err
	}
	if _, err := os.Create(filename); err != nil {
		return false, err
	}
	return true, nil
}
