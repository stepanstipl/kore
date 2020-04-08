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
	"strings"
)

func AppendMultiError(e1, e2 error) error {
	if e1 == nil {
		return e2
	}

	if e2 == nil {
		return e1
	}

	me1, ok := e1.(*multiError)
	if !ok {
		me1 = &multiError{Errors: []error{e1}}
	}
	me1.Errors = append(me1.Errors, e2)

	return me1
}

type multiError struct {
	Errors []error
}

func (m *multiError) Error() string {
	messages := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, ", ")
}
