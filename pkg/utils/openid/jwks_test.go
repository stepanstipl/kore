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

package openid

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertJWKSToPublicKey(t *testing.T) {
	key := `
	{
      "kid": "-99hWRwlsL2-5wXRzMeycoofzeRGaop0rfV68qMBpic",
      "kty": "RSA",
      "alg": "RS256",
      "use": "sig",
      "n": "kGWHUxaBsDfRZUel675nETzcg6mUhmq5MPoQ9uORlGgh8CK9PkopKDCgg0-HVPBMJGrfH4oVnYv0omPMtxigE6Q2mpxl-krwUT9jZdGomaqvS6_3lUy15oh-9r7nOv6CeupU61XsEPpUFQhRnMt7ZW2-sQnDNj8NOK_fhExcg4p6KL4wnHEk0r2AuR4KxzHOzUM80UTWSccM8lQ_3TLhCmRk8IgcmSHVd-PoMXGwf9Me5RkzbtbAdmSsfgkVxiQv9q3b2K1aLN3B2twvQku89RnDwe6Q7fyDWfj3uzvKO63K_RZ-IPrCt-WHrm-hB2NYKAtXkY4WIe9yMsb_8r4ejQ",
      "e": "AQAB",
      "x5t": "B4v8XETQ1nKJYK8UEk52w8A1mb8",
      "x5t#S256": "QKkhqLFYfBthA-YSPg1c23Qp6qe8TjKKdc8W9u5sS7g"
    }`
	values := make(map[string]string)

	require.NoError(t, json.NewDecoder(strings.NewReader(key)).Decode(&values))

	pkey, err := ConvertJWKSToPublicKey(values["n"], values["e"])
	require.NoError(t, err)
	require.NotNil(t, pkey)
}
