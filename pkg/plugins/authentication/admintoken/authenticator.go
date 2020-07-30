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

package admintoken

import (
	"context"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type authImpl struct {
	kore.Interface
	// config is the internal config
	config Config
}

// New returns a new header based identity provider
func New(h kore.Interface, config Config) (identity.Plugin, error) {
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

	id, found, err := o.GetUserIdentity(ctx, "admin", kore.WithAuthMethod("admintoken"))
	if err != nil || !found {
		return nil, false
	}

	return id, true
}

// Name returns the plugin name
func (o *authImpl) Name() string {
	return "admin-token"
}
