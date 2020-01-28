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
			Object:  &orgv1.TeamMemberList{},
			Version: orgv1.GroupVersion.String(),
			Kind:    "TeamMembershipList",
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
					"apiVersion": "org.hub.appvia.io/v1",
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
