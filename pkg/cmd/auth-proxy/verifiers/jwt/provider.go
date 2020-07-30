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

package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"
	"github.com/appvia/kore/pkg/utils"

	djwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type provider struct {
	Options
	// PublicKey is the key used to sign
	PublicKey *rsa.PublicKey
}

// New creates and returns an oidc provider
func New(options Options) (verifiers.Interface, error) {
	if options.Signer == nil {
		return nil, errors.New("no signer defined")
	}

	key, err := djwt.ParseRSAPublicKeyFromPEM(options.Signer)
	if err != nil {
		return nil, err
	}

	return &provider{Options: options, PublicKey: key}, nil
}

// Admit checks the token is valid
func (o *provider) Admit(request *http.Request) (bool, error) {
	// @step: extract the token from the request
	bearer, found := utils.GetBearerToken(request.Header.Get("Authorization"))
	if !found {
		return false, errors.New("no authorization token")
	}
	log.Debug("checking against the openid verifier")

	c := make(djwt.MapClaims)

	// @step: parse and extract the identity
	token, err := djwt.ParseWithClaims(bearer, &c, func(token *djwt.Token) (interface{}, error) {
		return o.PublicKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return false, nil
		}
	default:
		return false, nil
	}

	claims := utils.NewClaims(c)

	// @step: check the audience
	if aud, found := claims.GetAudience(); !found {
		log.Warn("no audience in the presented token")

		return false, nil
	} else if aud != "kubernetes" {
		log.Warn("invalid audience presented in the token")

		return false, nil
	}

	// @step: check the username
	username, found := claims.GetUserClaim("preferred_username", "email")
	if !found {
		return false, errors.New("no username found in the identity token")
	}

	// @step: ensure no impersonation is passed through by clearing all headers
	for name := range request.Header {
		if strings.HasPrefix(name, "Impersonate") {
			log.Warn("request has an impersonation header, denying the request")

			return false, nil
		}
	}

	request.Header.Set("Impersonate-User", username)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.ImpersonationToken))

	log.WithField("user", username).Debug("successfully authenticated user")

	return true, nil
}
