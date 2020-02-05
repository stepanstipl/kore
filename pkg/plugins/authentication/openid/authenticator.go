/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
		"discovery-url": config.DiscoveryURL,
		"user-claim":    config.UserClaims,
	}).Info("initializing the openid authentication plugin")

	// @step: grab an openid verifier
	discovery, err := openid.New(openid.Config{
		ClientID:     config.ClientID,
		DiscoveryURL: config.DiscoveryURL,
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
