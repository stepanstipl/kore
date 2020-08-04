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

package utils

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/appvia/kore/pkg/utils/openid"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims is used as a helper to JWT claims
type Claims struct {
	claims jwt.MapClaims
}

// NewClaims returns a claims
func NewClaims(claims jwt.MapClaims) *Claims {
	return &Claims{claims: claims}
}

// NewClaimsFromToken creates a claims from a openid.IDToken
func NewClaimsFromToken(token openid.IDToken) (*Claims, error) {
	c := jwt.MapClaims{}
	if err := token.Claims(&c); err != nil {
		return nil, err
	}

	return NewClaims(c), nil
}

// NewClaimsFromRawBytesToken returns a claims by parsing a raw token
func NewClaimsFromRawBytesToken(token []byte) (*Claims, error) {
	return NewClaimsFromRawToken(string(token))
}

// NewClaimsFromRawToken returns a claims by parsing a raw token
func NewClaimsFromRawToken(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return NewClaims(token.Claims.(jwt.MapClaims)), nil
}

// Sign is used to sign and mint a token
func (c *Claims) Sign(rsa []byte) ([]byte, error) {
	if !c.HasEmail() {
		return nil, errors.New("email required")
	}
	signer, err := jwt.ParseRSAPrivateKeyFromPEM(rsa)
	if err != nil {
		return nil, err
	}
	c.claims["alg"] = jwt.SigningMethodRS256.Alg()
	c.claims["iat"] = time.Now().UTC().Unix()

	minted := jwt.NewWithClaims(jwt.SigningMethodRS256, c.claims)
	token, err := minted.SignedString(signer)
	if err != nil {
		return nil, err
	}

	return []byte(token), nil
}

// GetUserClaim returns the username claim - defaults to 'name'
func (c *Claims) GetUserClaim(claims ...string) (string, bool) {
	for _, x := range claims {
		if name, found := c.GetString(x); found {
			return name, true
		}
	}

	return "", false
}

// GetAudience returns the aud of the jwt
func (c *Claims) GetAudience() (string, bool) {
	return c.GetString("aud")
}

// GetIssuer returns the iss
func (c *Claims) GetIssuer() (string, bool) {
	return c.GetString("iss")
}

// GetEmail returns the email claim
func (c *Claims) GetEmail() (string, bool) {
	return c.GetString("email")
}

// GetEmailVerified returns if the email is verified
func (c *Claims) GetEmailVerified() (bool, bool) {
	return c.GetBool("email_verified")
}

// HasEmail checks if the email exists
func (c *Claims) HasEmail() bool {
	_, found := c.GetEmail()

	return found
}

// HasExpired indicates the token has expired
func (c *Claims) HasExpired() bool {
	exp, found := c.GetExpiry()
	if !found {
		return false
	}

	return exp.Before(time.Now().UTC())
}

// GetExpiry returns the expiry of the jwt
func (c *Claims) GetExpiry() (time.Time, bool) {
	expiry, found := c.GetFloat64("exp")
	if !found {
		return time.Time{}, false
	}

	sec, dec := math.Modf(expiry)
	return time.Unix(int64(sec), int64(dec*(1e9))), true
}

// GetBool returns the boolean
func (c *Claims) GetBool(key string) (bool, bool) {
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

// GetFloat64 returns the float64 if found in the claims
func (c *Claims) GetFloat64(key string) (float64, bool) {
	v, found := c.claims[key]
	if !found {
		return 0, false
	}

	value, ok := v.(float64)

	if !ok {
		return 0, false
	}

	return value, true
}

// GetStringClaims trys to look for claims in token
func (c *Claims) GetStringClaims(keys ...string) (string, bool) {
	for _, name := range keys {
		if v, found := c.GetString(name); found {
			return v, true
		}
	}

	return "", false
}

// GetStringSlice returns a slice of string if found
func (c *Claims) GetStringSlice(key string) ([]string, bool) {
	v, found := c.claims[key]
	if !found {
		return []string{}, false
	}
	values, ok := v.([]string)
	if !ok {
		return []string{}, false
	}

	return values, true
}

// GetString returns the string from the claims
func (c *Claims) GetString(key string) (string, bool) {
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

// String returns the token itself
func (c *Claims) String() string {
	encoded, err := json.MarshalIndent(c.claims, "", "    ")
	if err != nil {
		return ""
	}

	return string(encoded)
}
