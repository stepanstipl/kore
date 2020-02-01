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

package openid

// Config is the configuration for the openid provider
type Config struct {
	// ClientID is the openid client id
	ClientID string
	// DiscoveryURL is the discovery URL
	DiscoveryURL string `json:"discovery-url,omitempty"`
	// SkipTLSVerify skips the TLS for the IDP
	SkipTLSVerify bool `json:"skip-tls-verify,omitempty"`
	// UserClaims is the claim fields which specifies the username
	UserClaims []string `json:"user-claim,omitempty"`
}
