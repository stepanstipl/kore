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

import "strings"

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
