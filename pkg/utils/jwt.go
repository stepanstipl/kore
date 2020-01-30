/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

package utils

import (
	"github.com/coreos/go-oidc"
	jwt "github.com/dgrijalva/jwt-go"
)

// Claims is used as a helper to JWT claims
type Claims struct {
	claims jwt.MapClaims
}

// NewClaims returns a claims
func NewClaims(claims jwt.MapClaims) Claims {
	return Claims{claims: claims}
}

func NewClaimsFromToken(token *oidc.IDToken) (Claims, error) {
	c := jwt.MapClaims{}
	if err := token.Claims(&c); err != nil {
		return Claims{}, err
	}

	return NewClaims(c), nil
}

// GetUserClaim returns the username claim - defaults to 'name'
func (c Claims) GetUserClaim(claims ...string) (string, bool) {
	if len(claims) <= 0 {
		return c.GetString("name")
	}
	for _, x := range claims {
		if name, found := c.GetString(x); found {
			return name, true
		}
	}

	return "", false
}

// GetEmail returns the email claim
func (c Claims) GetEmail() (string, bool) {
	return c.GetString("email")
}

// GetEmailVerified returns if the email is verified
func (c Claims) GetEmailVerified() (bool, bool) {
	return c.GetBool("email_verified")
}

// GetBool returns the boolean
func (c Claims) GetBool(key string) (bool, bool) {
	v, found := c.claims[key]
	if !found {
		return false, false
	}

	value, ok := v.(bool)
	if !ok {
		return false, false
	}

	return value, true
}

// GetString returns the string from the claims
func (c Claims) GetString(key string) (string, bool) {
	v, found := c.claims[key]
	if !found {
		return "", false
	}

	value, ok := v.(string)
	if !ok {
		return "", false
	}

	return value, true
}
