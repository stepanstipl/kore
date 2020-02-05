/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package kore

import (
	"errors"
)

// IsValid checks the options
func (c Config) IsValid() error {
	if c.DEX.EnabledDex && c.AdminPass == "" {
		return errors.New("you must set the admin password for dex")
	}
	if c.CertificateAuthority == "" {
		return errors.New("no certificate authority has been defined")
	}
	if c.CertificateAuthorityKey == "" {
		return errors.New("no certificate authority key")
	}

	return nil
}

// HasOpenID checks if openid is setup
func (c Config) HasOpenID() bool {
	return c.DiscoveryURL != ""
}

// HasCertificateAuthorityKey check if we have a key
func (c Config) HasCertificateAuthorityKey() bool {
	return c.CertificateAuthorityKey != ""
}

// HasCertificateAuthority checks if we have a CA
func (c Config) HasCertificateAuthority() bool {
	return c.CertificateAuthority != ""
}

// HasHMAC checks if the kore has a encoding token
func (c Config) HasHMAC() bool {
	return c.HMAC != ""
}
