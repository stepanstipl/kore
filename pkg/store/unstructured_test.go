/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func makeTestObject() metav1.Object {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "test",
			Labels: map[string]string{
				"hello": "world",
			},
		},
	}
}

func makeTestObjects() []metav1.Object {
	objects := make([]metav1.Object, 20)
	for i := 0; i < len(objects); i++ {
		name := fmt.Sprintf("test%d", i)
		objects[i] = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name": name,
				},
			},
		}
	}

	return objects
}

func TestObjectsToListOK(t *testing.T) {
	objects := makeTestObjects()

	err := ObjectsToList(&corev1.NamespaceList{}, objects)
	assert.NoError(t, err)
}

func TestObjectsToListUnstructuredBad(t *testing.T) {
	objects := makeTestObjects()

	err := ObjectsToList(&unstructured.Unstructured{}, objects)
	assert.Error(t, err)
}

func TestObjectsToListUnstructuredOK(t *testing.T) {
	objects := makeTestObjects()

	err := ObjectsToList(&unstructured.UnstructuredList{}, objects)
	assert.NoError(t, err)
}

func TestObjectToTypeOK(t *testing.T) {
	object := makeTestObject()

	err := ObjectToType(&corev1.Pod{}, object)
	require.NoError(t, err)
	assert.Equal(t, "pod", object.GetName())
	assert.Equal(t, "test", object.GetNamespace())
	assert.Equal(t, "world", object.GetLabels()["hello"])
}

func TestObjectToTypeUnstructuredOK(t *testing.T) {
	object := makeTestObject()

	err := ObjectToType(&unstructured.Unstructured{}, object)
	require.NoError(t, err)
	assert.Equal(t, "pod", object.GetName())
	assert.Equal(t, "test", object.GetNamespace())
}

func TestObjectToTypeBad(t *testing.T) {
	object := makeTestObject()
	err := ObjectToType(&corev1.Namespace{}, object)
	require.Error(t, err)
	assert.Equal(t, "invalid type, expected: v1.Namespace, needs: v1.Pod", err.Error())
}
