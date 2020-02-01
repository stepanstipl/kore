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

package identity

import (
	"context"
	"crypto/x509"
	"net/http"

	"github.com/appvia/kore/pkg/hub/authentication"
)

// Requestor is the interface for a request
type Requestor interface {
	// ClientCertficate is a potential client cert
	ClientCertficate() *x509.Certificate
	// Headers is a map of http header
	Headers() http.Header
}

// Plugin provides the interface for a authentication plugin, the purpose
// of which is to take an incoming bearer token and expand into a hub
// user identity
type Plugin interface {
	// Admit is a handler which passed
	Admit(context.Context, Requestor) (authentication.Identity, bool)
	// Name provides a name for the plugin
	Name() string
}
