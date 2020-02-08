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
