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

package server

import (
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/plugins/authentication/admintoken"
	"github.com/appvia/kore/pkg/plugins/authentication/basicauth"
	"github.com/appvia/kore/pkg/plugins/authentication/headers"
	"github.com/appvia/kore/pkg/plugins/authentication/openid"

	log "github.com/sirupsen/logrus"
)

// makeAuthenticators is responsible for configuration any authentication plugins
// @QUESTION: should we just move the configuration of theses into flags configured
// on the init()?
func makeAuthenticators(hubcc hub.Interface, config Config) error {
	if len(config.Hub.Authenticators) <= 0 {
		log.Warn("no authentication plugins have configured")

		return nil
	}

	// @step: we need to create any authentication plugins
	for _, x := range config.Hub.Authenticators {
		logger := log.WithFields(log.Fields{
			"plugin": x,
		})

		plugin, err := func() (identity.Plugin, error) {
			switch x {
			case "admintoken":
				return admintoken.New(hubcc, admintoken.Config{
					Token: config.Hub.AdminToken,
				})
			case "basicauth":
				return basicauth.New(hubcc)
			case "identity":
				return headers.New(hubcc)
			case "openid":
				return openid.New(hubcc, openid.Config{
					ClientID:     config.Hub.ClientID,
					DiscoveryURL: config.Hub.DiscoveryURL,
					UserClaims:   config.Hub.UserClaims,
				})
			default:
				return nil, errors.New("unknown plugin")
			}
		}()
		if err != nil {
			logger.WithError(err).Info("failed to register the authentication plugin")

			return err
		}

		// @step: register the plugin with the api server
		if err := identity.Register(plugin); err != nil {
			return fmt.Errorf("failed to register plugin %s, error: %s", x, err)
		}

		logger.Info("successfully registered the authentication plugin")
	}

	return nil
}
