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

	// @step: we need to verify the token
	token, err := o.verifier.Verify(ctx, bearer)
	if err != nil {
		return nil, false
	}

	claims, err := utils.NewClaimsFromToken(token)
	if err != nil {
		return nil, false
	}

	// @step: extract the user claims from the token
	username, found := claims.GetUserClaim(o.config.UserClaims...)
	if !found {
		log.Warn("issued token does not contain the username claim")

		return nil, false
	}

	// @step: find an sso identity with this username
	identity, found, err := o.GetUserIdentityByProvider(ctx, username, kore.IdentitySSO)
	if err != nil {
		log.WithError(err).Error("trying to find sso user in kore")

		return nil, false
	}
	if !found {
		return nil, false
	}

	// @step: retrieve the user identity
	user, found, err := o.GetUserIdentity(ctx, identity.User.Username, kore.WithAuthMethod("sso"))
	if err != nil {
		return nil, false
	}
	if !found {
		log.WithField(
			"username", identity.User.Username).
			Warn("sso identity was found but the user was not")

		return nil, false
	}

	return user, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "openid"
}
