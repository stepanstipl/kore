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
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/openid"
	oidc "github.com/appvia/kore/pkg/utils/openid"

	log "github.com/sirupsen/logrus"
)

// Options are options for the provider
type Options struct {
	// ClientID is the audience
	ClientID string
	// DiscoveryURL is the openid endpoint
	DiscoveryURL string
	// Token is the impersonation token
	Token string
	// UserClaims are the provider claims to use
	UserClaims []string
}

type provider struct {
	Options
	provider oidc.Authenticator
}

// New creates and returns an oidc provider
func New(options Options) (verifiers.Interface, error) {
	p, err := oidc.New(openid.Config{
		ClientID:          options.ClientID,
		ServerURL:         options.DiscoveryURL,
		SkipClientIDCheck: true,
	})
	if err != nil {
		return nil, err
	}

	if err := p.Run(context.Background()); err != nil {
		return nil, err
	}

	return &provider{Options: options, provider: p}, nil
}

// Admit checks the token is valid
func (o *provider) Admit(request *http.Request) (bool, error) {
	// @step: extract the token from the request
	bearer, found := utils.GetBearerToken(request.Header.Get("Authorization"))
	if !found {
		return false, errors.New("no authorization token")
	}
	log.Debug("checking against the openid verifier")

	// @step: parse and extract the identity
	id, err := o.provider.Verify(request.Context(), bearer)
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
	claims, err := utils.NewClaimsFromToken(id)
	if err != nil {
		return false, err
	}

	user, found := claims.GetUserClaim(o.UserClaims...)
	if !found {
		return false, errors.New("no username found in the identity token")
	}
	request.Header.Set("Impersonate-User", user)

	// @step: extract the group if request
	for _, x := range o.UserClaims {
		groups, found := claims.GetStringSlice(x)
		if found {
			for _, name := range groups {
				request.Header.Set("Impersonate-Group", name)
			}
		}
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.Token))

	log.WithField("user", user).Debug("successfully authenticated user")

	return true, nil
}
