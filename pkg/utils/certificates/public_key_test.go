/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package certificates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testKey = `-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAxz4coi0vBGONySMHH5T9q6cEEUZegCMIBhZixR0gUqBrrf079Rrg
MHsj3YaM0JLecr00QDs9emyWPdj1Oder76naJQcyFdX4qRqJeqPVnZZdzg7GSKNt
pl6DuOj+mCSxLf4U1DS12PG/ojM4keMUUxNBufgaaItXI1WqgkW7QAHb83gYARmL
nj+8/yiw/omRtpXWiYbM6M6nY/QQ5cEgZHRb+A0MsQGyw91DiQKnA09ljIRzvl+d
GuPn1vEhYuTPvvntLa8BdSiawgHzz7DG47U9/PwvjyLHEcLq5E2NhBnrJ2TGXdIb
4hAo4noChEGx+Yfnxm0SDKmlPv4u6LOyewIDAQAB
-----END PUBLIC KEY-----
`
	testPem = `-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAxz4coi0vBGONySMHH5T9q6cEEUZegCMIBhZixR0gUqBrrf079Rrg
MHsj3YaM0JLecr00QDs9emyWPdj1Oder76naJQcyFdX4qRqJeqPVnZZdzg7GSKNt
pl6DuOj+mCSxLf4U1DS12PG/ojM4keMUUxNBufgaaItXI1WqgkW7QAHb83gYARmL
nj+8/yiw/omRtpXWiYbM6M6nY/QQ5cEgZHRb+A0MsQGyw91DiQKnA09ljIRzvl+d
GuPn1vEhYuTPvvntLa8BdSiawgHzz7DG47U9/PwvjyLHEcLq5E2NhBnrJ2TGXdIb
4hAo4noChEGx+Yfnxm0SDKmlPv4u6LOyewIDAQAB
-----END PUBLIC KEY-----
`
)

func TestDecodePublicRSAKey(t *testing.T) {
	pk, err := DecodePublicRSAKey(strings.NewReader(testKey))
	require.NoError(t, err)
	require.NotNil(t, pk)
}

func TestEncodePublicKeyToPEM(t *testing.T) {
	pk, err := DecodePublicRSAKey(strings.NewReader(testKey))
	require.NoError(t, err)
	require.NotNil(t, pk)

	pem, err := EncodePublicKeyToPEM(pk)
	require.NoError(t, err)
	require.NotEmpty(t, pem)
	require.Equal(t, testPem, string(pem))
}
