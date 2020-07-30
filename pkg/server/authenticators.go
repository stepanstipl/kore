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

package server

import (
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/plugins/authentication/admintoken"
	"github.com/appvia/kore/pkg/plugins/authentication/basicauth"
	"github.com/appvia/kore/pkg/plugins/authentication/localjwt"
	"github.com/appvia/kore/pkg/plugins/authentication/openid"

	log "github.com/sirupsen/logrus"
)

// makeAuthenticators is responsible for configuration any authentication plugins
// @QUESTION: should we just move the configuration of theses into flags configured
// on the init()?
func makeAuthenticators(hubcc kore.Interface, config Config) error {
	if len(config.Kore.Authenticators) <= 0 {
		log.Warn("no authentication plugins have configured")

		return nil
	}

	// @step: we need to create any authentication plugins
	for _, x := range config.Kore.Authenticators {
		logger := log.WithFields(log.Fields{
			"plugin": x,
		})

		plugin, err := func() (identity.Plugin, error) {
			switch x {
			case "admintoken":
				return admintoken.New(hubcc, admintoken.Config{
					Token: config.Kore.AdminToken,
				})
			case "basicauth":
				return basicauth.New(hubcc)
			case "openid":
				return openid.New(hubcc, openid.Config{
					ClientID:   config.Kore.IDPClientID,
					ServerURL:  config.Kore.IDPServerURL,
					UserClaims: config.Kore.IDPUserClaims,
				})
			case "localjwt":
				return localjwt.New(hubcc, localjwt.Config{
					PublicKey:  config.Kore.LocalJWTPublicKey,
					UserClaims: config.Kore.IDPUserClaims,
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
