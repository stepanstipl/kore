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

package authproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeVerifier struct {
	allowed bool
}

func (f *fakeVerifier) Admit(req *http.Request) (bool, error) {
	if !f.allowed {
		return false, nil
	}

	return true, nil
}

func makeFakeVerifier(allowed bool) verifiers.Interface {
	return &fakeVerifier{allowed: allowed}
}

func makeFakeAuthProxy(t *testing.T, config *Config, v []verifiers.Interface) (*authImpl, *httptest.Server, *httptest.Server) {
	p, err := New(config, v)
	require.NoError(t, err)
	require.NotNil(t, p)

	s := httptest.NewServer(p.(*authImpl).handler)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Upstream", "true")
	}))

	config.UpstreamURL = upstream.URL

	err = p.(*authImpl).MakeRouter()
	require.NoError(t, err)

	return p.(*authImpl), s, upstream
}

func TestNew(t *testing.T) {
	_, s, p := makeFakeAuthProxy(t, &Config{}, []verifiers.Interface{makeFakeVerifier(false)})
	defer func() {
		s.Close()
		p.Close()
	}()
}

func TestServerAdmit(t *testing.T) {
	config := &Config{
		AllowedIPs: []string{"127.0.0.1/8"},
	}
	//allow := makeFakeVerifier(true)
	disallow := makeFakeVerifier(false)

	cases := []struct {
		Expected int
	}{
		{Expected: http.StatusForbidden},
		//{Expected: http.StatusOK},
	}
	_, proxy, upstream := makeFakeAuthProxy(t, config, []verifiers.Interface{disallow})
	defer func() {
		proxy.Close()
		upstream.Close()
	}()

	for _, c := range cases {
		req, err := http.NewRequest(http.MethodGet, proxy.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, c.Expected, resp.StatusCode)
	}
}
