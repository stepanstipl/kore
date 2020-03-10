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

package identity

import (
	"context"
	"crypto/x509"
	"net/http"

	"github.com/appvia/kore/pkg/kore/authentication"
)

// Requestor is the interface for a request
type Requestor interface {
	// ClientCertficate is a potential client cert
	ClientCertficate() *x509.Certificate
	// Headers is a map of http header
	Headers() http.Header
}

// Plugin provides the interface for a authentication plugin, the purpose
// of which is to take an incoming bearer token and expand into a kore
// user identity
type Plugin interface {
	// Admit is a handler which passed
	Admit(context.Context, Requestor) (authentication.Identity, bool)
	// Name provides a name for the plugin
	Name() string
}
