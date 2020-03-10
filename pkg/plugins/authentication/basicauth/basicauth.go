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
	"encoding/base64"
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
