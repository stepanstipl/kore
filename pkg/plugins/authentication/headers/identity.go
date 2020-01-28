/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

package headers

import (
	"context"

	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/hub/authentication"
	plugin "github.com/appvia/kore/pkg/apiserver/plugins/identity"
)

type hdrImpl struct {
	hub.Interface
}

// New returns a new header based identity provider
func New(h hub.Interface) (plugin.Plugin, error) {
	return &hdrImpl{Interface: h}, nil
}

// Admit is called to authenticate the inbound request
func (h hdrImpl) Admit(ctx context.Context, req plugin.Requestor) (authentication.Identity, bool) {
	// @step: grab the identity header from the request
	username := req.Headers().Get("X-Identity")
	if username == "" {
		return nil, false
	}
	identity, found, err := hub.Client.GetUserIdentity(ctx, username)
	if err != nil || !found {
		return nil, false
	}

	return identity, true
}

// Name returns the plugin name
func (h hdrImpl) Name() string {
	return "identity"
}
