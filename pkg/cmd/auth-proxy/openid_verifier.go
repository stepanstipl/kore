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

package authproxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/utils"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
)

type openidImpl struct {
	token    string
	claims   []string
	verifier *oidc.IDTokenVerifier
}

// NewOpenIDAuth creates an openid provider
func NewOpenIDAuth(clientID, endpoint, token string, claims []string) (AuthProvider, error) {
	options := &oidc.Config{
		ClientID:          clientID,
		SkipClientIDCheck: true,
		SkipExpiryCheck:   false,
	}
	log.WithField(
		"idp-server-url", endpoint,
	).Info("using the IDP server to verify the requests")

	if endpoint == "" {
		return nil, errors.New("no endpoint defined")
	}
	if token == "" {
		return nil, errors.New("no token for impersonation")
	}

	provider, err := oidc.NewProvider(context.Background(), token)
	if err != nil {
		log.WithError(err).Error("trying to retrieve provider details")

		return nil, err
	}

	return &openidImpl{
		claims:   claims,
		token:    token,
		verifier: provider.Verifier(options),
	}, nil
}

// Admit checks the token is valid
func (o *openidImpl) Admit(request *http.Request) (bool, error) {
	// @step: extract the token from the request
	bearer, found := utils.GetBearerToken(request.Header.Get("Authorization"))
	if !found {
		return false, errors.New("no authorization token")
	}

	// @step: parse and extract the identity
	idToken, err := o.verifier.Verify(request.Context(), bearer)
	if err != nil {
		return false, err
	}

	// @step: ensure no impersonation is passed through by clearing all headers
	for name := range request.Header {
		if strings.HasPrefix(name, "Impersonate") {
			request.Header.Del(name)
		}
	}

	// @step: extract the username if any
	claims, err := utils.NewClaimsFromToken(idToken)
	if err != nil {
		return false, err
	}

	user, found := claims.GetUserClaim(o.claims...)
	if !found {
		return false, errors.New("no username found in the identity token")
	}
	request.Header.Set("Impersonate-User", user)

	// @step: extract the group if request
	for _, x := range o.claims {
		groups, found := claims.GetStringSlice(x)
		if found {
			for _, name := range groups {
				request.Header.Set("Impersonate-Group", name)
			}
		}
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.token))

	return true, nil
}
