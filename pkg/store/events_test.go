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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAddNamespace(t *testing.T) {
	s := newTestStore(t)
	doneCh := make(chan metav1.Object)

	s.WatchResource("v1/namespaces")
	err := s.AddEventListener(&Listener{
		EventHandlers: &EventHandlerFuncs{
			CreatedHandlerFunc: func(object metav1.Object) {
				doneCh <- object
			},
		},
		Resources: []string{"*"},
	})
	require.NoError(t, err)

	ns, err := s.client.CoreV1().Namespaces().Create(
		context.Background(),
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
		},
		metav1.CreateOptions{},
	)

	require.NotNil(t, ns)
	require.NoError(t, err)

	select {
	case e := <-doneCh:
		require.NotNil(t, e)
		assert.Equal(t, "test", e.GetName())
	case <-time.After(time.Millisecond * 100):
		t.Error("failed to update the store on namespace change")
	}
}

func TestDeletedNamespace(t *testing.T) {
	doneCh := make(chan metav1.Object)

	s := newTestStore(t)
	s.WatchResource("v1/namespaces")
	err := s.AddEventListener(&Listener{
		EventHandlers: &EventHandlerFuncs{
			DeletedHandlerFunc: func(object metav1.Object) {
				doneCh <- object
			},
		},
	})
	require.NoError(t, err)

	ns, err := s.client.CoreV1().Namespaces().Create(
		context.Background(),
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test"}},
		metav1.CreateOptions{},
	)
	require.NotNil(t, ns)
	require.NoError(t, err)

	err = s.client.CoreV1().Namespaces().Delete(context.Background(), "test", metav1.DeleteOptions{})
	require.NoError(t, err)

	select {
	case e := <-doneCh:
		require.NotNil(t, e)
		assert.Equal(t, "test", e.GetName())
	case <-time.After(time.Millisecond * 100):
		t.Error("failed to delete the store on namespace change")
	}
}
