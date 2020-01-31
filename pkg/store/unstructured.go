/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package store

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// listItemsField is the name of the items field
	listItemsField = "Items"
)

var (
	// convertor is the unstructured converter client
	converter = runtime.DefaultUnstructuredConverter
)

// ObjectToType is responsible for converting the cached metav1.Objects
// into the correct types
func ObjectToType(holder runtime.Object, object metav1.Object) error {
	// @step: reflect the placeholder to get the type
	v := reflect.ValueOf(holder)
	// get a struct reference
	st := v.Elem()

	// @step: check if the placeholder is an typed or unstructured list - we assume no
	unstruct := false
	if st.Type() == reflect.TypeOf(unstructured.Unstructured{}) {
		unstruct = true
	}

	if unstruct {
		converted, err := converter.ToUnstructured(object)
		if err != nil {
			return fmt.Errorf("failed to convert metav1.Object to Unstructured: %s", err)
		}
		// we get a reference to the field
		uobj := st.FieldByName("Object")
		// we create an new unstructured.Unstructured
		uobj.Set(reflect.ValueOf(converted))

		return nil
	}

	// @guard: add a check to ensure we asking for and getting the same types
	expected := v.Elem().Type()
	needs := reflect.ValueOf(object).Elem().Type()
	if expected != needs {
		return fmt.Errorf("invalid type, expected: %s, needs: %s", expected, needs)
	}

	// we check the object supported the interface
	dc, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object: %T does not support the runtime.Object interface", object)
	}
	// we deepcopy the object so the nothing upstream effects the cache
	copied := dc.DeepCopyObject()
	// update the placeholder
	v.Elem().Set(reflect.Indirect(reflect.ValueOf(copied)))

	return nil
}

// ObjectsToList is responsible for deep copying and convertng the objects
// to a list of typed or unstructured items
func ObjectsToList(holder runtime.Object, objects []metav1.Object) error {
	if objects == nil {
		return nil
	}

	// @step: reflect the placeholder to get the type
	v := reflect.ValueOf(holder)

	// @step: ensure the value has a items field else we throw it away
	// as an invalid type
	field, found := v.Type().Elem().FieldByName(listItemsField)
	if !found {
		return fmt.Errorf("invalid runtime.Object, no %s field", listItemsField)
	}
	// get a struct reference
	st := v.Elem()
	// @step: get a reference to the items field
	items := st.FieldByName(listItemsField)

	// @step: check if the placeholder is an typed or unstructured list - we assume no
	unstruct := false
	if field.Type == reflect.TypeOf([]unstructured.Unstructured{}) {
		unstruct = true
	}

	if unstruct {
		obj := map[string]interface{}{
			"kind":       "List",
			"apiVersion": "v1",
		}
		st.FieldByName("Object").Set(reflect.ValueOf(obj))
	} else {
		st.FieldByName("Kind").Set(reflect.ValueOf("List"))
		st.FieldByName("APIVersion").Set(reflect.ValueOf("v1"))
	}

	// @step: we iterate the objects and inject into the reflected field
	for _, x := range objects {
		switch unstruct {
		case true:
			// we need to convert the object to a unstructured type
			converted, err := converter.ToUnstructured(x)
			if err != nil {
				return fmt.Errorf("failed to convert metav1.Object to Unstructured: %s", err)
			}
			// we create an new unstructured.Unstructured
			item := unstructured.Unstructured{
				Object: converted,
			}
			// we append to the list
			items.Set(reflect.Append(items, reflect.ValueOf(item)))
		default:
			// we check the object supported the interface
			dc, ok := x.(runtime.Object)
			if !ok {
				return fmt.Errorf("object: %T does not support the runtime.Object interface", x)
			}
			// we deepcopy the object so the nothing upstream effects the cache
			copied := dc.DeepCopyObject()
			// we append the object to the reflected field
			items.Set(reflect.Append(items, reflect.Indirect(reflect.ValueOf(copied))))
		}
	}

	return nil
}
