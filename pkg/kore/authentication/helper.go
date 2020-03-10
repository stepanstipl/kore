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

package authentication

import "context"

// GetIdentity returns the identity from the context if any
func GetIdentity(ctx context.Context) (Identity, bool) {
	v := ctx.Value(ContextKey{})
	if v == nil {
		return nil, false
	}

	return v.(Identity), true
}

// MustGetIdentity returns the identity from the context if any
func MustGetIdentity(ctx context.Context) Identity {
	v := ctx.Value(ContextKey{})
	if v == nil {
		panic("no user identity found in the context")
	}

	return v.(Identity)
}
