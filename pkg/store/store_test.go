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

package store

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/appvia/kore/pkg/schema"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
	crfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newTestStore(t *testing.T) *storeImpl {
	s, err := New(fake.NewSimpleClientset(), crfake.NewFakeClientWithScheme(schema.GetScheme()))
	require.NotNil(t, s)
	require.NoError(t, err)

	return s.(*storeImpl)
}

func TestNew(t *testing.T) {
	s, err := New(fake.NewSimpleClientset(), crfake.NewFakeClient())
	assert.NotNil(t, s)
	assert.NoError(t, err)
}

func TestStoreList(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 10; i++ {
		node := &unstructured.Unstructured{}
		node.SetName(fmt.Sprintf("nothing%d", i))
		require.NoError(t, s.Kind("nothing").Set(node.GetName(), node))
	}

	for i := 0; i < 10; i++ {
		node := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "node.io",
				"kind":       "node",
				"metadata": map[string]interface{}{
					"labels": map[string]string{},
				},
			},
		}
		node.SetName(fmt.Sprintf("node%d", i))
		require.NoError(t, s.Kind("Node").Set(node.GetName(), node))
	}

	items, err := s.Kind("Node").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.NotEmpty(t, items)
	assert.Equal(t, 10, len(items))
}

func TestAPIGroup(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 10; i++ {
		o := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "node.io",
				"metadata": map[string]interface{}{
					"name":   fmt.Sprintf("node%d", i),
					"labels": map[string]string{},
				},
			},
		}
		require.NoError(t, s.APIVersion("node.io").Kind("node").Set(o.GetName(), o))
	}
	items, err := s.Kind("node").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.NotEmpty(t, items)
	assert.Equal(t, 10, len(items))

	items, err = s.APIVersion("node.io").Kind("node").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.NotEmpty(t, items)
	assert.Equal(t, 10, len(items))
}

func TestStoreGetByLabels(t *testing.T) {
	s := newTestStore(t)

	o := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]string{
				"name": "test",
			},
		},
	}
	oa := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]string{"name": "test1"}},
	}
	s.APIVersion("kore.io").Kind("Workspace").Namespace("kore").Label("team.io", "dev").Set("test", o)
	s.APIVersion("kore.io").Kind("Workspace").Namespace("kore").Label("team.io", "apps").Set("test1", oa)

	items, err := s.APIVersion("kore.io").Kind("Workspace").Namespace("kore").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.Equal(t, 2, len(items))

	items, err = s.APIVersion("kore.io").Kind("Workspace").Namespace("kore").Label("team.io", "dev").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.Equal(t, 1, len(items))

	items, err = s.APIVersion("kore.io").Kind("Workspace").Namespace("kore").Label("team.io", "apps").List()
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.Equal(t, 1, len(items))
}

func TestStoreGet(t *testing.T) {
	s := newTestStore(t)

	node := &unstructured.Unstructured{}
	node.SetName("node1")
	node.SetKind("Node")
	s.Kind("nodes").Set("node1", node)

	v, err := s.Kind("nodes").Get("node1")
	require.NoError(t, err)
	require.NotNil(t, v)
	assert.Equal(t, "node1", v.GetName())
	ty := v.(metav1.Type)
	assert.Equal(t, "Node", ty.GetKind())
}

func TestStoreSet(t *testing.T) {
	s := newTestStore(t)

	node := &unstructured.Unstructured{}
	node.SetName("node1")
	s.Kind("nodes").Set("node1", node)

	v, err := s.Kind("nodes").Get("node1")
	require.NoError(t, err)
	require.NotNil(t, v)
	assert.Equal(t, "node1", v.GetName())
}

func TestStoreActions(t *testing.T) {
	cs := []struct {
		Actions func(s Store)
		Checks  func(s Store)
	}{
		{
			Actions: func(s Store) {
				s.Kind("nodes").Set("node1", &unstructured.Unstructured{})
				s.Kind("nodes").Set("node2", &unstructured.Unstructured{})
			},
			Checks: func(s Store) {
				found, err := s.Kind("nodes").Has("node2")
				assert.True(t, found)
				assert.NoError(t, err)
				items, err := s.Kind("nodes").List()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(items))
			},
		},
		{
			Actions: func(s Store) {
				s.Kind("nodes").Set("node1", &unstructured.Unstructured{})
				s.Kind("nodes").Set("node2", &unstructured.Unstructured{})
				require.NoError(t, s.Kind("nodes").Delete("node1"))
			},
			Checks: func(s Store) {
				found, err := s.Kind("nodes").Has("node1")
				assert.False(t, found)
				assert.NoError(t, err)
				items, err := s.Kind("nodes").List()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(items))
			},
		},
		{
			Actions: func(s Store) {
				s.Kind("namespaces").Set("default", &unstructured.Unstructured{})
			},
			Checks: func(s Store) {
				found, err := s.Kind("namespaces").Has("default")
				assert.True(t, found)
				assert.NoError(t, err)
			},
		},
		{
			Actions: func(s Store) {
				s.Kind("namespaces").Set("default", &unstructured.Unstructured{})
				for _, x := range []string{"test0", "test1"} {
					s.Namespace("default").Kind("pods").Set(x, &unstructured.Unstructured{})
				}
			},
			Checks: func(s Store) {
				v, err := s.Kind("namespaces").Has("default")
				require.NoError(t, err)
				require.NotNil(t, v)
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(items))
			},
		},
		{
			Actions: func(s Store) {
				s.Namespace("default").Kind("pods").Set("test0", &unstructured.Unstructured{})
				s.Namespace("default").Kind("pods").Set("test1", &unstructured.Unstructured{})
				require.NoError(t, s.Namespace("default").Kind("pods").Delete("test0"))
			},
			Checks: func(s Store) {
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(items))
			},
		},
		{
			Actions: func(s Store) {
				require.NoError(t, s.Kind("namespaces").Delete("default"))
				for _, x := range []string{"test0", "test1"} {
					require.NoError(t, s.Namespace("default").Kind("pods").Delete(x))
				}
			},
			Checks: func(s Store) {
				found, err := s.Kind("namespaces").Has("default")
				require.False(t, found)
				require.NoError(t, err)
				items, err := s.Namespace("default").Kind("pods").List()
				assert.NoError(t, err)
				assert.Equal(t, 0, len(items))
			},
		},
	}
	for _, c := range cs {
		s := newTestStore(t)
		c.Actions(s)
		c.Checks(s)
	}
}
