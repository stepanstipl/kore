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

package localjwt

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type authImpl struct {
	kore.Interface
	// config is the configuration
	config Config
}

// New returns an jwt authenticator
func New(h kore.Interface, config Config) (identity.Plugin, error) {
	// @step: verify the configuration
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	log.Info("initializing the jwt authentication plugin")

	return &authImpl{Interface: h, config: config}, nil
}

// Admit is called to authenticate the inbound request
func (o *authImpl) Admit(ctx context.Context, req identity.Requestor) (authentication.Identity, bool) {

	// @step: verify the authorization token
	bearer, found := utils.GetBearerToken(req.Headers().Get("Authorization"))
	if !found {
		return nil, false
	}

	id, err := func() (authentication.Identity, error) {
		// @step: validate the token
		claims := struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			jwt.StandardClaims
		}{}
		_, err := jwt.ParseWithClaims(bearer, &claims, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			pubKeyBytes, err := base64.StdEncoding.DecodeString(o.config.PublicKey)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse public key from config: %v", err)
			}
			pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse public key from config: %v", err)
			}
			return pubKey, nil
		})
		if err != nil {
			return nil, err
		}

		username := claims.Username
		if username == "" {
			return nil, errors.New("issued token does not contain the username claim")
		}

		id, found, err := o.GetUserIdentity(ctx, username, kore.WithAuthMethod("localjwt"))
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, errors.New("user not found in the kore")
		}

		return id, nil
	}()
	if err != nil {
		return nil, false
	}

	return id, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "jwt"
}
