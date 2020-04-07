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

	"github.com/appvia/kore/pkg/kore"

	"github.com/urfave/cli/v2"
)

// GetDeleteAdminCommand returns the create member command
func GetDeleteAdminCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "admin",
		Aliases: []string{"admins"},
		Usage:   "Deletes an administrator from kore",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "The username of the user you wish to add to the team",
				Required: true,
			},
		},

		Action: func(ctx *cli.Context) error {
			username := ctx.String("user")

			err := NewRequest().
				WithConfig(config).
				WithEndpoint("/teams/{team}/members/{user}").
				PathParameter("team", true).
				PathParameter("user", true).
				WithInject("team", kore.HubAdminTeam).
				WithInject("user", username).
				Delete()
			if err != nil {
				return err
			}

			fmt.Printf("User %q has been removed as an administrator from kore\n", username)

			return nil
		},
	}
}
