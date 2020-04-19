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

package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFakeUpstream(t *testing.T) *httptest.Server {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Hit", "true")
			w.WriteHeader(http.StatusOK)
		}))

	require.NotNil(t, s)

	return s
}

func TestNewNoEndpoint(t *testing.T) {
	v, err := New(Options{})
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestInvalidEndpoint(t *testing.T) {
	v, err := New(Options{Endpoint: "$$$"})
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestReverseProxy(t *testing.T) {
	s := makeFakeUpstream(t)
	defer s.Close()
	v, err := New(Options{Endpoint: s.URL})
	require.NoError(t, err)
	require.NotNil(t, v)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	v.Serve(nil).ServeHTTP(resp, req)

	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	assert.Equal(t, resp.Header().Get("Hit"), "true")
}
