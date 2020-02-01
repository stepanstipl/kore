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
