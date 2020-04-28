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

package kubernetes

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/appvia/kore/pkg/utils"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	// ErrFieldNotFound indicates the field was not found in the resource
	ErrFieldNotFound = errors.New("field not found")
)

// GetRuntimeField is used to extract a type from a runtime.Object
func GetRuntimeField(object runtime.Object, path string, out interface{}) error {
	u := &unstructured.Unstructured{}

	if !utils.IsEqualType(object, &unstructured.Unstructured{}) {
		// @step: apply the cluster configuration to the component
		us, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
		if err != nil {
			return err
		}
		u.Object = us
	} else {
		u = object.(*unstructured.Unstructured)
	}

	value, found, err := unstructured.NestedFieldCopy(u.Object, strings.Split(path, ".")...)
	if err != nil {
		return err
	}
	if !found {
		return ErrFieldNotFound
	}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(value); err != nil {
		return err
	}

	return json.Unmarshal(b.Bytes(), out)
}
