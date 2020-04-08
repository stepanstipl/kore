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
	"github.com/appvia/kore/pkg/kore"

	"github.com/urfave/cli/v2"
)

// GetCreateAdminCommand returns the create member command
func GetCreateAdminCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "admin",
		Aliases: []string{"admins"},
		Usage:   "Creates a new administrator in kore",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "The username of the user you wish to add to the team",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "invite",
				Aliases:  []string{"i"},
				Usage:    "If the user doesn't exist and the invite flag is set, the invite url will be automatically generated.",
				Required: false,
			},
		},

		Action: func(ctx *cli.Context) error {
			invite := ctx.Bool("invite")
			username := ctx.String("user")

			return CreateMember(config, kore.HubAdminTeam, username, invite)
		},
	}
}
