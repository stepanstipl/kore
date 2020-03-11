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

package schema

import (
	"io/ioutil"
	"testing"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestGetSchema(t *testing.T) {
	assert.NotNil(t, GetScheme())
}

func TestIsVersionedOK(t *testing.T) {
	assert.True(t, IsVersioned(&orgv1.User{}))
	assert.True(t, IsVersioned(&orgv1.TeamMember{}))

	unstruct := &unstructured.Unstructured{}
	unstruct.SetAPIVersion(orgv1.GroupVersion.String())
	unstruct.SetKind("User")
	assert.True(t, IsVersioned(unstruct))
}

func TestIsVersionedBad(t *testing.T) {
	unstruct := &unstructured.Unstructured{}
	unstruct.SetAPIVersion("notthere/v1")
	unstruct.SetKind("Nokind")

	assert.False(t, IsVersioned(unstruct))
}

func TestGetGroupKindVersionBad(t *testing.T) {
	o := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "unknown.io/v1",
			"kind":       "Nothing",
		},
	}
	kind, registered, err := GetGroupKindVersion(o)
	assert.NoError(t, err)
	assert.False(t, registered)
	assert.NotNil(t, kind)
}

func TestGetGroupKindVersionOK(t *testing.T) {
	cases := []struct {
		Kind    string
		Object  runtime.Object
		Version string
	}{
		{
			Object:  &orgv1.User{},
			Version: orgv1.GroupVersion.String(),
			Kind:    "User",
		},
		{
			Object:  &orgv1.Team{},
			Version: orgv1.GroupVersion.String(),
			Kind:    "Team",
		},
		{
			Object:  &configv1.Allocation{},
			Version: configv1.GroupVersion.String(),
			Kind:    "Allocation",
		},
		{
			Object:  &configv1.Plan{},
			Version: configv1.GroupVersion.String(),
			Kind:    "Plan",
		},
		{
			Object: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "org.kore.appvia.io/v1",
					"kind":       "User",
				},
			},
			Version: orgv1.GroupVersion.String(),
			Kind:    "User",
		},
	}
	for _, c := range cases {
		kind, versioned, err := GetGroupKindVersion(c.Object)
		require.NoError(t, err)
		require.True(t, versioned)
		assert.Equal(t, c.Version, kind.GroupVersion().String())
		assert.Equal(t, c.Kind, kind.GroupKind().Kind)
	}
}
