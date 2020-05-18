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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
)

// GetRuntimeKind returns the object kind
func GetRuntimeKind(o runtime.Object) string {
	return o.GetObjectKind().GroupVersionKind().Kind
}

// GetRuntimeName returns the runtime name
func GetRuntimeName(o runtime.Object) string {
	meta, err := GetMeta(o)
	if err != nil {
		return ""
	}

	return meta.Name
}

// GetRuntimeSelfLink returns the self link
func GetRuntimeSelfLink(o runtime.Object) (string, error) {
	meta, err := GetMeta(o)
	if err != nil {
		return "", err
	}
	gvk := o.GetObjectKind().GroupVersionKind()

	return fmt.Sprintf("%s/%s/%s/%s", gvk.Group, gvk.Version, strings.ToLower(gvk.Kind), meta.Name), nil
}

// MustGetRuntimeSelfLink returns the self link
func MustGetRuntimeSelfLink(o runtime.Object) string {
	gvk := o.GetObjectKind().GroupVersionKind()

	return fmt.Sprintf("%s/%s/%s/%s", gvk.Group, gvk.Version, strings.ToLower(gvk.Kind), GetRuntimeName(o))
}
