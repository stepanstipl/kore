/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package openid

import (
	"context"

	"github.com/coreos/go-oidc"
)

// Authenticator provides a openid if agent
type Authenticator interface {
	// Run start the discovery grab
	Run(context.Context) error
	// Verify is called to verify a token
	Verify(context.Context, string) (*oidc.IDToken, error)
}

// Config is the configuration for the service
type Config struct {
	// ClientID is the client id
	ClientID string
	// DiscoveryURL is the openid discovery url
	DiscoveryURL string
	// SkipIDCheck indicates we skip checking the issuer
	SkipClientIDCheck bool
}
