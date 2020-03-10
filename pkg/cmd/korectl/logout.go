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

package korectl

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// GetLogoutCommand is used to login to the api server
func GetLogoutCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:      "logout",
		Usage:     "Deletes the login credentials from the current or selected profile",
		ArgsUsage: "[PROFILE]",
		Action: func(ctx *cli.Context) error {
			var profileName = config.CurrentProfile
			if ctx.Args().Present() {
				profileName = ctx.Args().First()
			}

			profile, ok := config.Profiles[profileName]
			if !ok {
				return fmt.Errorf("%q profile does not exist", profileName)
			}
			authInfo, ok := config.AuthInfos[profile.AuthInfo]
			if !ok {
				return fmt.Errorf("%q user does not exist in the korectl configuration", profile.AuthInfo)
			}

			authInfo.OIDC = &OIDC{
				ClientID:     authInfo.OIDC.ClientID,
				ClientSecret: authInfo.OIDC.ClientSecret,
				AuthorizeURL: authInfo.OIDC.AuthorizeURL,
			}

			if err := config.Update(); err != nil {
				return err
			}

			fmt.Println("Logout successful.")

			return nil
		},
	}
}
