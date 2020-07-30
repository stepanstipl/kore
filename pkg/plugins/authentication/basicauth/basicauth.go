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

package basicauth

import (
	"context"
	"strings"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"
)

type authImpl struct {
	kore.Interface
}

// New returns a new header based identity provider
func New(h kore.Interface) (identity.Plugin, error) {
	return &authImpl{Interface: h}, nil
}

// Admit is called to authenticate the inbound request
func (o *authImpl) Admit(ctx context.Context, req identity.Requestor) (authentication.Identity, bool) {
	// @step: verify the authorization token
	username, password, found := utils.GetBasicAuthFromHeader(req.Headers().Get("Authorization"))
	if !found {
		return nil, false
	}

	id, found, err := o.GetUserIdentityByProvider(ctx, username, kore.IdentityBasicAuth)
	if err != nil || !found {
		return nil, false
	}

	// @TODO we rework this when i do the rbac piece as it will involve moving the
	// authentication piece around alot - as we'd to wrap in some form of encryption service
	current := id.ProviderToken
	switch {
	case strings.HasPrefix(current, "md5:"):
		hash := strings.TrimPrefix(current, "md5:")
		if utils.HashString(password) != hash {
			return nil, false
		}
	default:
		if current != password {
			return nil, false
		}
	}

	user, found, err := o.GetUserIdentity(ctx, username, kore.WithAuthMethod("basicauth"))
	if err != nil || !found {
		return nil, false
	}

	return user, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "basicauth"
}
