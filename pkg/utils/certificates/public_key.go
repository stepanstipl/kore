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
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

// DecodePublicRSAKey decodes the pem encoded public key
func DecodePublicRSAKey(reader io.Reader) (*rsa.PublicKey, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, errors.New("no encoded pem found in content")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func EncodePublicKeyToPEM(pubkey *rsa.PublicKey) ([]byte, error) {
	asn1Bytes, err := asn1.Marshal(*pubkey)
	if err != nil {
		return []byte{}, err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	b := &bytes.Buffer{}
	if err = pem.Encode(b, pemkey); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}
