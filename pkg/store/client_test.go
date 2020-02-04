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
	"context"
	"fmt"
	"testing"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func makeTestPod() *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testpod",
			Namespace: "test",
			Labels: map[string]string{
				"hello": "world",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "test",
		},
	}
}

func TestClient(t *testing.T) {
	s := newTestStore(t)
	c := s.Client()
	assert.NotNil(t, c)
}

func TestClientCreateNoValue(t *testing.T) {
	s := newTestStore(t)
	c := s.Client()

	require.NotNil(t, c)
	require.NoError(t, s.WatchResource("v1/pods"))
	require.Error(t, c.Create(context.TODO()))
}

func TestClientCreateOK(t *testing.T) {
	s := newTestStore(t)
	pod := makeTestPod()
	c := s.Client()

	require.NotNil(t, c)
	require.NoError(t, s.WatchResource("v1/pods"))
	require.NoError(t, c.Create(context.TODO(), CreateOptions.From(pod)))

	apod := &corev1.Pod{}
	require.NoError(t, c.Get(context.TODO(),
		GetOptions.InNamespace(pod.Namespace),
		GetOptions.InTo(apod),
		GetOptions.WithName(pod.Name),
	))

	assert.Equal(t, pod, apod)
}

func TestClientCreateWithCacheOK(t *testing.T) {
	s := newTestStore(t)
	pod := makeTestPod()
	c := s.Client()

	require.NotNil(t, c)
	require.NoError(t, c.Create(context.TODO(), CreateOptions.From(pod)))

	apod := &corev1.Pod{}
	require.NoError(t, c.Get(context.TODO(),
		GetOptions.InNamespace(pod.Namespace),
		GetOptions.InTo(apod),
		GetOptions.WithCache(true),
		GetOptions.WithName(pod.Name),
	))
	assert.Equal(t, pod, apod)

	pb := &dto.Metric{}
	cacheHitCounter.Write(pb)
	assert.Equal(t, float64(1), pb.GetCounter().GetValue())
}

func TestClientListUnstructuredOK(t *testing.T) {
	s := newTestStore(t)

	// @step: populate the cache and api with users
	for i := 0; i < 3; i++ {
		username := fmt.Sprintf("name%d", i)
		err := s.Client().Create(context.TODO(),
			CreateOptions.From(&orgv1.User{
				TypeMeta: metav1.TypeMeta{
					APIVersion: orgv1.SchemeGroupVersion.String(),
					Kind:       "User",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      username,
					Namespace: "kore",
				},
			}),
		)
		require.NoError(t, err)
	}

	// @step: ask for them in an unstructured form
	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   orgv1.GroupVersion.Group,
		Kind:    "UserList",
		Version: orgv1.GroupVersion.Version,
	})

	require.NoError(t, s.Client().List(context.TODO(),
		ListOptions.InNamespace("kore"),
		ListOptions.InTo(list),
		ListOptions.WithCache(true),
	))
	assert.Equal(t, len(list.Items), 3)

	nocache := &unstructured.UnstructuredList{}
	nocache.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   orgv1.GroupVersion.Group,
		Kind:    "UserList",
		Version: orgv1.GroupVersion.Version,
	})

	/*
		require.NoError(t, s.Client().List(context.TODO(),
			ListOptions.InNamespace("kore"),
			ListOptions.InTo(nocache),
			ListOptions.WithCache(false),
		))
		assert.Equal(t, len(nocache.Items), 3)
	*/
}

func TestClientListTypedOK(t *testing.T) {
	s := newTestStore(t)

	for i := 0; i < 10; i++ {
		username := fmt.Sprintf("name%d", i)
		err := s.Client().Create(context.TODO(),
			CreateOptions.From(&orgv1.User{
				TypeMeta: metav1.TypeMeta{
					APIVersion: orgv1.SchemeGroupVersion.String(),
					Kind:       "User",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      username,
					Namespace: "kore",
				},
			}),
		)
		require.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("team%d", i)
		err := s.Client().Create(context.TODO(),
			CreateOptions.From(&orgv1.Team{
				TypeMeta: metav1.TypeMeta{
					APIVersion: orgv1.SchemeGroupVersion.String(),
					Kind:       "Team",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "kore",
				},
			}),
		)
		require.NoError(t, err)
	}

	list := &orgv1.UserList{}
	require.NoError(t, s.Client().List(context.TODO(),
		ListOptions.InNamespace("kore"),
		ListOptions.WithCache(true),
		ListOptions.InTo(list),
	))
	assert.Equal(t, len(list.Items), 10)
}
