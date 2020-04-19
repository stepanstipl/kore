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

import (
	"errors"
)

// IsValid checks the options are valid
func (o *Options) IsValid() error {
	if o.DiscoveryURL == "" {
		return errors.New("no server url")
	}
	if o.ClientID == "" {
		return errors.New("no client id")
	}
	if len(o.UserClaims) == 0 {
		return errors.New("no user claims")
	}

	return nil
}
