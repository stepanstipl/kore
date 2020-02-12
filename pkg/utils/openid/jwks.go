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
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"math/big"
)

// ConvertJWKSToPublicKey converts the keys from the jwks url to a public key
func ConvertJWKSToPublicKey(e, n string) (*rsa.PublicKey, error) {
	// @step: decode the base64 of the n
	decN, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}

	pn := big.NewInt(0)
	pn.SetBytes(decN)

	decE, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}

	var eBytes []byte
	switch len(decE) < 8 {
	case true:
		eBytes = make([]byte, 8-len(decE), 8)
		eBytes = append(eBytes, decE...)
	default:
		eBytes = decE
	}

	eReader := bytes.NewReader(eBytes)
	var pe uint64

	if err = binary.Read(eReader, binary.BigEndian, &pe); err != nil {
		return nil, err
	}

	return &rsa.PublicKey{N: pn, E: int(pe)}, nil
}
