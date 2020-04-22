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
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReflectedItem struct {
	Name    string
	Age     int
	Pointer *string
	Array   []string
}

func setterFunc(field string, value interface{}) func() (string, interface{}) {
	return func() (string, interface{}) {
		return field, value
	}
}

func TestIsEqual(t *testing.T) {
	assert.False(t, IsEqualType(ReflectedItem{}, nil))
	assert.False(t, IsEqualType(nil, &ReflectedItem{}))
	assert.True(t, IsEqualType(&ReflectedItem{}, &ReflectedItem{}))
	assert.True(t, IsEqualType(&ReflectedItem{}, ReflectedItem{}))
	assert.True(t, IsEqualType(ReflectedItem{}, &ReflectedItem{}))
	assert.True(t, IsEqualType(nil, nil))
}

func TestSetReflectedFieldOK(t *testing.T) {
	message := "test"

	cases := []struct {
		Item     ReflectedItem
		Expected ReflectedItem
		Setters  []func() (string, interface{})
	}{
		{
			Item:     ReflectedItem{},
			Expected: ReflectedItem{},
		},
		{
			Item: ReflectedItem{},
			Setters: []func() (string, interface{}){
				setterFunc("Name", "Hello"),
			},
			Expected: ReflectedItem{Name: "Hello"},
		},
		{
			Item: ReflectedItem{},
			Setters: []func() (string, interface{}){
				setterFunc("Name", "Hello"),
				setterFunc("Age", 30),
			},
			Expected: ReflectedItem{Name: "Hello", Age: 30},
		},
		{
			Item: ReflectedItem{},
			Setters: []func() (string, interface{}){
				setterFunc("Pointer", &message),
			},
			Expected: ReflectedItem{Pointer: &message},
		},
		{
			Item: ReflectedItem{},
			Setters: []func() (string, interface{}){
				setterFunc("Array", []string{"a", "b"}),
			},
			Expected: ReflectedItem{Array: []string{"a", "b"}},
		},
	}

	for _, c := range cases {
		for _, setter := range c.Setters {
			field, value := setter()
			SetReflectedField(field, value, &c.Item)
		}
		assert.Equal(t, c.Expected, c.Item)
	}
}

func TestHasReflectedField(t *testing.T) {
	assert.False(t, HasReflectField("NotThere", &ReflectedItem{}))
	assert.False(t, HasReflectField("NotThere", ReflectedItem{}))
	assert.True(t, HasReflectField("Name", &ReflectedItem{}))
	assert.True(t, HasReflectField("Name", ReflectedItem{}))
}
