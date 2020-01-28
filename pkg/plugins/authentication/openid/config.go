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

package openid

import (
	"errors"
	"fmt"
	"net/url"
)

// IsValid checks the configuration is valid
func (c Config) IsValid() error {
	if c.DiscoveryURL == "" {
		return errors.New("no discovery url configured")
	}
	if c.ClientID == "" {
		return errors.New("no client id configured")
	}
	if _, err := url.Parse(c.DiscoveryURL); err != nil {
		return fmt.Errorf("invalid discovery url: %s", err)
	}
	if len(c.UserClaims) <= 0 {
		c.UserClaims = append(c.UserClaims, []string{"preferred_username"}...)
	}

	return nil
}
