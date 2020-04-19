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

package health

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDefaultNext() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Next", "OK")
		w.WriteHeader(http.StatusOK)
	})
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New())
}

func TestReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	resp := httptest.NewRecorder()

	New().Serve(makeDefaultNext()).ServeHTTP(resp, req)
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	require.NotNil(t, resp.Body)

	content, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte("OK\n"), content)
	assert.Empty(t, resp.Header().Get("Next"))
}

func TestNotReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	New().Serve(makeDefaultNext()).ServeHTTP(resp, req)
	assert.Equal(t, resp.Result().StatusCode, http.StatusOK)
	require.NotNil(t, resp.Body)
	assert.Equal(t, resp.Header().Get("Next"), "OK")
}
