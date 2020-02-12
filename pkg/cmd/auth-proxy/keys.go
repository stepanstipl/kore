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

package authproxy

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/appvia/kore/pkg/utils/certificates"

	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
)

type staticKeys struct {
	verifier *rsa.PublicKey
}

// newStaticKeySet returns a custom key set
func newStaticKeySet(path string) (oidc.KeySet, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	verifier, err := certificates.DecodePublicRSAKey(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return &staticKeys{verifier: verifier}, nil
}

// VerifySignature is used to verify the token
func (s *staticKeys) VerifySignature(_ context.Context, bearer string) ([]byte, error) {
	_, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		return s.verifier, nil
	})
	if err != nil {
		return []byte{}, err
	}

	parts := strings.Split(bearer, ".")

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %v", err)
	}

	return payload, nil
}
