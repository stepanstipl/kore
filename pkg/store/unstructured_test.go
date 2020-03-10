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
