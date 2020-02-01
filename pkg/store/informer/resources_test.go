/*
Copyright 2018 Appvia Ltd <info@appvia.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package informer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceNames(t *testing.T) {
	e := ResourceNames()
	assert.NotEmpty(t, e)
}

func TestResourceVersions(t *testing.T) {
	e := ResourceVersions()
	require.NotNil(t, e)
	assert.NotEmpty(t, e)
}

func TestToSchema(t *testing.T) {
	cs := []struct {
		Resource string
		Expected schema.GroupVersionResource
	}{
		{
			Resource: "v1/pods",
			Expected: corev1.SchemeGroupVersion.WithResource("pods"),
		},
		{
			Resource: "core/v1/pods",
			Expected: schema.GroupVersionResource{
				Group:    "core",
				Resource: "pods",
				Version:  "v1",
			},
		},
	}
	for i, c := range cs {
		s, err := ToSchema(c.Resource)
		require.NoErrorf(t, err, "case %d, did not expect error: %s", i, err)
		assert.Equal(t, c.Expected, s)
	}
}

func TestToSchemaBad(t *testing.T) {
	for _, k := range []string{"bad", "bad-v1", "bad-d1/r1/r1/r3"} {
		_, err := ToSchema(k)
		assert.Error(t, err)
	}
}

func TestResourceVersionsName(t *testing.T) {
	e := ResourceVersions()["v1/namespaces"]
	expected := corev1.SchemeGroupVersion.WithResource("namespaces")
	assert.Equal(t, expected, e)
}
