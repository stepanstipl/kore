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

package tokenreview

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func makeFakeVerifier(t *testing.T, items ...runtime.Object) (kubernetes.Interface, *kvImpl) {
	client := fake.NewSimpleClientset(items...)
	v, err := NewFromClient(client, Options{})
	require.NoError(t, err)
	require.NotNil(t, v)

	return client, v.(*kvImpl)
}

func TestNewFromClient(t *testing.T) {
	client := fake.NewSimpleClientset()
	v, err := NewFromClient(client, Options{})
	require.NoError(t, err)
	assert.NotNil(t, v)
}

func TestAdmitNoToken(t *testing.T) {
	_, v := makeFakeVerifier(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	permitted, err := v.Admit(req)

	require.NoError(t, err)
	require.False(t, permitted)
}

/*
func TestAdmitFailed(t *testing.T) {
	failure := &authentication.TokenReview{}
	_, v := makeFakeVerifier(t, failure)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid")

	permitted, err := v.Admit(req)

	require.NoError(t, err)
	require.False(t, permitted)
}
*/
