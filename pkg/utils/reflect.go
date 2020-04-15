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
	"errors"
	"reflect"
)

var (
	// ErrNoValue indicate no value was set
	ErrNoValue = errors.New("no value defined")
)

// SetReflectedField checks if the interface has the field
func SetReflectedField(name string, value, o interface{}) {
	var caller reflect.Value

	if reflect.ValueOf(o).Kind() == reflect.Ptr {
		caller = reflect.ValueOf(o).Elem()
	} else {
		caller = reflect.ValueOf(o)
	}

	field := caller.FieldByName(name)

	if !field.IsValid() || !field.CanSet() {
		return
	}
	field.Set(reflect.ValueOf(value))
}

// HasReflectField checks if a field exists in a interface
func HasReflectField(name string, o interface{}) bool {
	if reflect.ValueOf(o).Kind() == reflect.Ptr {
		return reflect.ValueOf(o).Elem().FieldByName(name).IsValid()
	}

	return reflect.ValueOf(o).FieldByName(name).IsValid()
}
