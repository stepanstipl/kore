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

package authproxy

import (
	"errors"
)

// IsValid checks the configuation of the proxy
func (c Config) IsValid() error {
	if c.ClientID == "" {
		return errors.New("no client id")
	}
	if c.DiscoveryURL == "" && c.SigningCA == "" {
		return errors.New("neither disovery-url or signing ca are not defined")
	}
	if len(c.UserClaims) <= 0 {
		return errors.New("user claims are empty")
	}
	if c.TLSCert != "" && c.TLSKey == "" {
		return errors.New("no tls private key")
	}
	if c.TLSKey != "" && c.TLSCert == "" {
		return errors.New("no tls certificate")
	}
	if c.UpstreamURL == "" {
		return errors.New("no upstream url")
	}

	return nil
}

// HasTLS checks if we have tls
func (c Config) HasTLS() bool {
	return c.TLSCert != "" && c.TLSKey != ""
}
