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
