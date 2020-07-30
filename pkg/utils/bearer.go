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
	"encoding/base64"
	"strings"
)

// GetBearerToken returns the bearer token from an authorization header
func GetBearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}

	items := strings.Split(header, " ")
	if len(items) != 2 {
		return "", false
	}

	if strings.ToLower(items[0]) != "bearer" {
		return "", false
	}

	return items[1], true
}

// GetBasicAuthToken is used to retrieve the basic authentication
func GetBasicAuthToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, " ")
	if len(items) != 2 {
		return "", false
	}

	if strings.ToLower(items[0]) != "basic" {
		return "", false
	}

	return items[1], true
}

// GetBasicAuthFromHeader returns the basic auth from authorization
func GetBasicAuthFromHeader(header string) (string, string, bool) {
	auth, found := GetBasicAuthToken(header)
	if !found {
		return "", "", false
	}

	payload, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", "", false
	}
	keypair := strings.SplitN(string(payload), ":", 2)
	if len(keypair) != 2 {
		return "", "", false
	}

	return keypair[0], keypair[1], true
}
