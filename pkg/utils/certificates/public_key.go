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
