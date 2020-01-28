/*
Copyright 2018 Rohith Jayawardene <gambol99@gmail.com>

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
	"context"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNew(t *testing.T) {
	client := fake.NewSimpleClientset()

	inf, err := New(&Config{
		Factories: []informers.SharedInformerFactory{informers.NewSharedInformerFactoryWithOptions(client, 0)},
		Resource:  "v1/namespaces",
	})
	require.NoError(t, err)
	require.NoError(t, inf.Stop())
}

func TestNewUnknownResource(t *testing.T) {
	inf, err := New(&Config{Resource: "unknown"})
	require.Error(t, err)
	assert.Nil(t, inf)
}

func TestMultipleInformers(t *testing.T) {
	client := fake.NewSimpleClientset()

	doneCh := make(chan struct{})
	factory := informers.NewSharedInformerFactoryWithOptions(client, 0)

	// @step: create the namespace informer
	podInf, err := New(&Config{
		Factories: []informers.SharedInformerFactory{factory},
		Resource:  "v1/namespaces",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
		},
	})
	require.NoError(t, err)
	require.NotNil(t, podInf)
	defer func() {
		require.NoError(t, podInf.Stop())
	}()

	nsInf, err := New(&Config{
		Factories: []informers.SharedInformerFactory{factory},
		Resource:  "v1/pods",
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "test_pod", object.GetName())
			doneCh <- struct{}{}
		},
	})
	require.NoError(t, err)
	require.NotNil(t, nsInf)
	defer func() {
		require.NoError(t, nsInf.Stop())
	}()

	// @step: add the namespace and pod
	_, err = client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NoError(t, err)

	_, err = client.CoreV1().Pods("default").Create(&v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test_pod"}})
	require.NoError(t, err)

	select {
	case <-time.After(1000 * time.Millisecond):
		t.Error("failed to recieve the done signal with time period")
	case <-doneCh:
	}
}

func TestInformerCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	doneCh := make(chan struct{})
	errorCh := make(chan error)
	client := fake.NewSimpleClientset()

	inf := newTestInformer(t, &Config{
		Resource:  "v1/namespaces",
		Factories: []informers.SharedInformerFactory{informers.NewSharedInformerFactoryWithOptions(client, 0)},
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
			doneCh <- struct{}{}
		},
		ErrorFunc: func(version schema.GroupVersionResource, err error) {
			errorCh <- err
		},
	})
	defer func() {
		assert.NoError(t, inf.Stop())
	}()

	// @step: add a namespace and check we get a update
	ns, err := client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NotNil(t, ns)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
		cancel()
	}
}

func TestInformerUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	var updated int
	doneCh := make(chan struct{})
	errorCh := make(chan error)
	client := fake.NewSimpleClientset()

	inf := newTestInformer(t, &Config{
		Resource:  "v1/namespaces",
		Factories: []informers.SharedInformerFactory{informers.NewSharedInformerFactoryWithOptions(client, 0)},
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			updated++
		},
		UpdateFunc: func(version schema.GroupVersionResource, before, after metav1.Object) {
			updated++
			if updated == 3 {
				require.NotNil(t, before)
				require.NotNil(t, after)
				assert.Equal(t, "default", before.GetName())
				annotations := after.GetAnnotations()
				require.NotNil(t, annotations)
				assert.Equal(t, "default", after.GetName())
				assert.Equal(t, "test", after.GetAnnotations()["test"])
				doneCh <- struct{}{}
			}
		},
	})
	defer inf.Stop()

	// @step: add a namespace and check we get a update
	ns, err := client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NotNil(t, ns)
	require.NoError(t, err)
	ns.SetAnnotations(map[string]string{"test": "test"})
	_, err = client.CoreV1().Namespaces().Update(ns)
	require.NoError(t, err)
	ns.SetAnnotations(map[string]string{"test": "test"})
	_, err = client.CoreV1().Namespaces().Update(ns)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func TestInformerDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	client := fake.NewSimpleClientset()

	doneCh := make(chan struct{})
	errorCh := make(chan error)
	inf := newTestInformer(t, &Config{
		Resource:  "v1/namespaces",
		Factories: []informers.SharedInformerFactory{informers.NewSharedInformerFactoryWithOptions(client, 0)},
		DeleteFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			require.NotNil(t, object)
			assert.Equal(t, "default", object.GetName())
			doneCh <- struct{}{}
		},
	})
	defer inf.Stop()

	// @step: add a namespace and check we get a update
	_, err := client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "default"}})
	require.NoError(t, err)
	require.NoError(t, client.CoreV1().Namespaces().Delete("default", &metav1.DeleteOptions{}))

	select {
	case <-ctx.Done():
		assert.NoError(t, ctx.Err())
	case err := <-errorCh:
		assert.NoError(t, err)
	case <-doneCh:
	}
}

func newTestInformer(t *testing.T, c *Config) Informer {
	inf, err := New(c)
	require.NoError(t, err)
	require.NotNil(t, inf)

	return inf
}
