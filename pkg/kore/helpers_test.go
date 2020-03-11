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

package kore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestEmptyUser(t *testing.T) {
	u := EmptyUser("test")
	require.NotNil(t, u)
	assert.Equal(t, "test", u.Name)
	assert.Equal(t, HubNamespace, u.Namespace)
}

func TestIsValidGVKOK(t *testing.T) {
	gvk := schema.GroupVersionKind{
		Group:   "something",
		Version: "v1",
		Kind:    "something",
	}
	assert.NoError(t, IsValidGVK(gvk))
}

func TestUnstructuredKind(t *testing.T) {
	kind := schema.GroupVersionKind{
		Group:   "things",
		Kind:    "Something",
		Version: "v1",
	}
	u := UnstructuredKind(kind)
	require.NotNil(t, u)
	assert.Equal(t, "things/v1", u.GetAPIVersion())
	assert.Equal(t, "Something", u.GetKind())
}

func TestIsValidGVKBad(t *testing.T) {
	cases := []struct {
		GVK schema.GroupVersionKind
	}{
		{GVK: schema.GroupVersionKind{}},
		{GVK: schema.GroupVersionKind{Group: "something"}},
		{GVK: schema.GroupVersionKind{Version: "v1"}},
		{GVK: schema.GroupVersionKind{Kind: "something"}},
		{GVK: schema.GroupVersionKind{Group: "a", Version: "v1"}},
		{GVK: schema.GroupVersionKind{Group: "a", Kind: "b"}},
	}
	for _, c := range cases {
		assert.Error(t, IsValidGVK(c.GVK))
	}
}

func TestLabel(t *testing.T) {
	assert.Equal(t, "kore.appvia.io/me", Label("me"))
}
