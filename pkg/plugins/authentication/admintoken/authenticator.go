/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package admintoken

import (
	"context"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type authImpl struct {
	hub.Interface
	// config is the internal config
	config Config
}

// New returns a new header based identity provider
func New(h hub.Interface, config Config) (identity.Plugin, error) {
	if config.Token == "" {
		config.Token = utils.Random(32)

		log.WithFields(log.Fields{
			"token": config.Token,
		}).Warn("no admin token has been defined, generate a ephermal one")
	}

	return &authImpl{Interface: h, config: config}, nil
}

// Admit is called to authenticate the inbound request
func (o *authImpl) Admit(ctx context.Context, req identity.Requestor) (authentication.Identity, bool) {
	// @step: verify the authorization token
	bearer, found := utils.GetBearerToken(req.Headers().Get("Authorization"))
	if !found {
		return nil, false
	}

	if bearer != o.config.Token {
		return nil, false
	}

	id, found, err := o.GetUserIdentity(ctx, "admin")
	if err != nil || !found {
		return nil, false
	}

	return id, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "admin-token"
}
