/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

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
