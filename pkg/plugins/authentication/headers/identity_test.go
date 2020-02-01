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
