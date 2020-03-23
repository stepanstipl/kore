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

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

// GetCreateTeamMemberCommand returns the create member command
func GetCreateTeamMemberCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Creates a new team member",

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
			team := ctx.String("team")
			username := ctx.String("user")

			if team == "" {
				return errTeamParameterMissing
			}

			// @step: check the team exist
			found, err := ResourceExists(config, "team", team)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("team %q does not exist", team)
			}

			// @step: check the user exist
			found, err = ResourceExists(config, "user", username)
			if err != nil {
				return err
			}
			if !found {
				if !invite {
					prompt := promptui.Prompt{
						Label:     "The user does not exist. Do you want to create an invite link",
						IsConfirm: true,
						Default:   "Y",
					}

					// Prompt will return an error if the input is not y/Y
					if _, err := prompt.Run(); err != nil {
						return nil
					}
				}

				var inviteURL string
				err := NewRequest().
					WithConfig(config).
					WithContext(ctx).
					WithEndpoint("/teams/{team}/invites/generate/{user}").
					PathParameter("team", true).
					WithInject("team", team).
					PathParameter("user", true).
					WithRuntimeObject(&inviteURL).
					Get()
				if err != nil {
					return err
				}
				fmt.Printf("Invite URL: %s\n", inviteURL)
				return nil
			}

			err = NewRequest().
				WithConfig(config).
				WithContext(ctx).
				WithEndpoint("/teams/{team}/members/{user}").
				PathParameter("team", true).
				PathParameter("user", true).
				Update()
			if err != nil {
				return err
			}
			fmt.Printf("%q has been successfully added to team %q\n", ctx.String("user"), ctx.String("team"))

			return nil
		},
	}
}
