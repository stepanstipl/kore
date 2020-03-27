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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// IsMissingKind checks if a runtime kind error
func IsMissingKind(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Object 'Kind' is missing in")
}

// GetOwnership returns an ownership from a object
func GetOwnership(object *unstructured.Unstructured, name string) corev1.Ownership {
	field, found, err := unstructured.NestedFieldNoCopy(object.Object, "spec", name)
	if err != nil || !found {
		return corev1.Ownership{}
	}

	owner, ok := field.(corev1.Ownership)
	if !ok {
		return corev1.Ownership{}
	}

	return owner
}

// InjectValuesIntoUnstructured injects the values into the spec
func InjectValuesIntoUnstructured(values map[string]interface{}, object *unstructured.Unstructured) {
	spec, found := object.Object["spec"].(map[string]interface{})
	if !found {
		spec = make(map[string]interface{})
		object.Object["spec"] = spec
	}

	for k, v := range values {
		spec[k] = v
	}
}

// InjectOwnershipIntoUnstructured is used to inject the oweership
func InjectOwnershipIntoUnstructured(name string, owner corev1.Ownership, object *unstructured.Unstructured) {
	spec, ok := object.Object["spec"].(map[string]interface{})
	if !ok {
		spec = make(map[string]interface{})
		object.Object["spec"] = spec
	}

	provider, ok := spec[name].(corev1.Ownership)
	if !ok {
		provider = corev1.Ownership{}
		spec[name] = corev1.Ownership{
			Group:     owner.Group,
			Kind:      owner.Kind,
			Version:   owner.Version,
			Name:      owner.Name,
			Namespace: owner.Namespace,
		}

		return
	}

	provider.Group = owner.Group
	provider.Kind = owner.Kind
	provider.Version = owner.Version
	provider.Name = owner.Name
	provider.Namespace = owner.Namespace
}
