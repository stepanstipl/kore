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
