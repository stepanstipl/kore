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
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultiCharSet is the default charset to use
	DefaultiCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomWithCharset returns a random string of x charset
func RandomWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

// Random returns a random string
func Random(length int) string {
	return RandomWithCharset(length, DefaultiCharSet)
}

// Rand returns a randon length of digits
func Rand(length int) (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	rand := strings.ReplaceAll(u.String(), "-", "")
	if length > len(rand) {
		length = len(rand)
	}

	return rand[:length], nil
}
