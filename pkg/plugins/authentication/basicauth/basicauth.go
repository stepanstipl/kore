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

package basicauth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/utils"
)

type authImpl struct {
	hub.Interface
}

// New returns a new header based identity provider
func New(h hub.Interface) (identity.Plugin, error) {
	return &authImpl{Interface: h}, nil
}

// Admit is called to authenticate the inbound request
func (o *authImpl) Admit(ctx context.Context, req identity.Requestor) (authentication.Identity, bool) {
	// @step: verify the authorization token
	basic, found := utils.GetBasicAuthToken(req.Headers().Get("Authorization"))
	if !found {
		return nil, false
	}

	payload, err := base64.StdEncoding.DecodeString(basic)
	if err != nil {
		return nil, false
	}
	keypair := strings.SplitN(string(payload), ":", 2)
	if len(keypair) != 2 {
		return nil, false
	}
	username := keypair[0]
	password := keypair[1]

	id, found, err := o.GetUserIdentityByProvider(ctx, username, "basicauth")
	if err != nil || !found {
		return nil, false
	}

	if id.ProviderToken != password {
		return nil, false
	}

	user, found, err := o.GetUserIdentity(ctx, username)
	if err != nil || !found {
		return nil, false
	}

	return user, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "basicauth"
}
