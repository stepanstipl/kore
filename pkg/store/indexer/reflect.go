/*
Copyright 2018 Appvia Ltd <info@appvia.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package indexer

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// buildReflectQuery is responsible for reflect the struct and building a struct from it
func buildReflectQuery(search interface{}) (string, error) {
	// @check we have a reference to a struct
	if reflect.ValueOf(search).Kind() != reflect.Ptr &&
		reflect.ValueOf(search).Elem().Kind() != reflect.Struct {
		return "", errors.New("search must be a type struct")
	}

	var terms []string

	v := reflect.ValueOf(search).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		termName := strings.ToLower((v.Type().Field(i).Name))
		termValue := ""

		var fieldValue = reflect.ValueOf(field.Interface())
		switch kind := field.Type().Kind(); kind {
		case reflect.Bool:
			termValue = fmt.Sprintf("%t", fieldValue.Bool())
		case reflect.String:
			termValue = field.String()
		case reflect.Int:
			termValue = fmt.Sprintf("%d", fieldValue.Int())
		case reflect.Map:
			// @check this is a map[string]string or nil
			if field.Type().String() != "map[string]string" {
				continue
			}
			for _, x := range fieldValue.MapKeys() {
				mv := fieldValue.MapIndex(x).String()
				if mv == "" {
					continue
				}
				terms = append(terms, fmt.Sprintf("+%s.%s:%s", termName, x.String(), mv))
			}
			continue
		}

		// @choice for numeric values we choose to ignore zero values
		// @choice for empty string we ignore
		if termValue == "0" || termValue == "" {
			continue
		}

		//fieldName := strings.ToLower(v.Type().Field(i).Name)
		terms = append(terms, fmt.Sprintf("+%s:%s", termName, termValue))
	}

	return strings.Join(terms, " "), nil
}
