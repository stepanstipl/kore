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

package headers

import (
	"context"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type fakeRequestor struct {
	header http.Header
	cert   *x509.Certificate
}

func (f fakeRequestor) Headers() http.Header {
	return f.header
}

func (f fakeRequestor) ClientCertficate() *x509.Certificate {
	return f.cert
}

func TestNew(t *testing.T) {
	p, err := New(nil)
	assert.NotNil(t, p)
	assert.NoError(t, err)
}

func TestAdmitBad(t *testing.T) {
	p, err := New(nil)
	require.NotNil(t, p)
	require.NoError(t, err)

	fake := &fakeRequestor{
		header: make(http.Header, 0),
	}

	id, found := p.Admit(context.TODO(), fake)
	assert.False(t, found)
	assert.Nil(t, id)
}

/*
func TestAdmitOK(t *testing.T) {
	p, err := New(nil)
	require.NotNil(t, p)
	require.NoError(t, err)

	fake := &fakeRequestor{
		header: map[string][]string{
			"X-Identity": {"test"},
		},
	}
	id, found := p.Admit(context.TODO(), fake)
	require.True(t, found)
	assert.Equal(t, "test", id.Username())
}
*/
