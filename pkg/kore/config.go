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

package kore

import (
	"errors"

	"github.com/appvia/kore/pkg/utils"
)

// IsValid checks the options
func (c Config) IsValid() error {
	if c.DEX.EnabledDex {
		if c.AdminPass == "" {
			return errors.New("you must set the admin password for dex")
		}
		if utils.Contains("offline", c.ClientScopes) {
			return errors.New("'offline' scope when using dex should be 'offline_access'")
		}
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
