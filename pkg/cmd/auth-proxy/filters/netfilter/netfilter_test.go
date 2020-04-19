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

package netfilter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	v, err := New(Options{})
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestInvalidCIDR(t *testing.T) {
	v, err := New(Options{Permitted: []string{"bad"}})
	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = New(Options{Permitted: []string{"1.1.1.1"}})
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestAdmit(t *testing.T) {
	v, err := New(Options{
		Permitted: []string{
			"127.0.0.1/8",
			"192.168.0.0/16",
			"1.1.1.1/32",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, v)

	next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Next", "OK")
	})

	cases := []struct {
		Address  string
		Expected int
	}{
		{Address: "127.0.0.1:2323", Expected: http.StatusOK},
		{Address: "117.0.0.1:2323", Expected: http.StatusForbidden},
		{Address: "1.1.1.1:2323", Expected: http.StatusOK},
		{Address: "1.1.1.2:2323", Expected: http.StatusForbidden},
		{Address: "bas:2323", Expected: http.StatusForbidden},
		{Address: "", Expected: http.StatusForbidden},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = c.Address

		resp := httptest.NewRecorder()
		v.Serve(next).ServeHTTP(resp, req)
		assert.Equal(t, c.Expected, resp.Result().StatusCode)

		switch c.Expected {
		case http.StatusOK:
			assert.Equal(t, "OK", resp.Header().Get("Next"))
		default:
			assert.Empty(t, resp.Header().Get("Next"))
		}
	}
}
