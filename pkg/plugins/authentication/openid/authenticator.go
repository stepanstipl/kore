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

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/openid"

	log "github.com/sirupsen/logrus"
)

type authImpl struct {
	kore.Interface
	// config is the configuration
	config Config
	// verifier is the openid verifer
	verifier openid.Authenticator
}

// New returns an openid authenticator
func New(h kore.Interface, config Config) (identity.Plugin, error) {
	// @step: verify the configuration
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"server-url": config.ServerURL,
		"user-claim": config.UserClaims,
	}).Info("initializing the openid authentication plugin")

	// @step: grab an openid verifier
	discovery, err := openid.New(openid.Config{
		ClientID:  config.ClientID,
		ServerURL: config.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	// @step: start the grab
	if err := discovery.Run(context.Background()); err != nil {
		return nil, err
	}

	return &authImpl{Interface: h, config: config, verifier: discovery}, nil
}

// Admit is called to authenticate the inbound request
func (o *authImpl) Admit(ctx context.Context, req identity.Requestor) (authentication.Identity, bool) {

	// @step: verify the authorization token
	bearer, found := utils.GetBearerToken(req.Headers().Get("Authorization"))
	if !found {
		return nil, false
	}

	id, err := func() (authentication.Identity, error) {
		// @step: valiidate the token
		token, err := o.verifier.Verify(ctx, bearer)
		if err != nil {
			return nil, err
		}

		claims, err := utils.NewClaimsFromToken(token)
		if err != nil {
			return nil, err
		}

		username, found := claims.GetUserClaim(o.config.UserClaims...)
		if !found {
			return nil, errors.New("issued token does not contain the username claim")
		}

		id, found, err := o.GetUserIdentity(ctx, username)
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
	return "openid"
}
