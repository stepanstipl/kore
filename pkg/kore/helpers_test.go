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
