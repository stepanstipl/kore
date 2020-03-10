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

package headers

import (
	"context"

	plugin "github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
)

type hdrImpl struct {
	kore.Interface
}

// New returns a new header based identity provider
func New(h kore.Interface) (plugin.Plugin, error) {
	return &hdrImpl{Interface: h}, nil
}

// Admit is called to authenticate the inbound request
func (h hdrImpl) Admit(ctx context.Context, req plugin.Requestor) (authentication.Identity, bool) {
	// @step: grab the identity header from the request
	username := req.Headers().Get("X-Identity")
	if username == "" {
		return nil, false
	}
	identity, found, err := kore.Client.GetUserIdentity(ctx, username)
	if err != nil || !found {
		return nil, false
	}

	return identity, true
}

// Name returns the plugin name
func (h hdrImpl) Name() string {
	return "identity"
}
