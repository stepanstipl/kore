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

package openid

// Config is the configuration for the openid provider
type Config struct {
	// ClientID is the openid client id
	ClientID string
	// ServerURL is the openid server URL
	ServerURL string `json:"server-url,omitempty"`
	// SkipTLSVerify skips the TLS for the IDP
	SkipTLSVerify bool `json:"skip-tls-verify,omitempty"`
	// UserClaims is the claim fields which specifies the username
	UserClaims []string `json:"user-claim,omitempty"`
}
