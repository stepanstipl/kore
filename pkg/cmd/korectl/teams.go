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
	"net/http"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/urfave/cli/v2"
)

// GetEditTeamCommand returns the edit team command
func GetEditTeamCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:      "team",
		Aliases:   []string{"teams"},
		Usage:     "Modifies a team",
		ArgsUsage: "TEAM",
		Action: func(ctx *cli.Context) error {
			teamID := ctx.Args().First()

			team := &orgv1.Team{}

			req := NewRequest().
				WithConfig(config).
				WithContext(ctx).
				PathParameter("id", true).
				WithInject("id", teamID).
				WithEndpoint("/teams/{id}").
				WithRuntimeObject(team)
			if err := req.Get(); err != nil {
				if reqErr, ok := err.(*RequestError); ok {
					if reqErr.statusCode == http.StatusNotFound {
						return fmt.Errorf("%q does not exist", teamID)
					}
				}
				return err
			}

			prompts := prompts{
				&prompt{id: "Description", value: &team.Spec.Description},
			}

			if err := prompts.collect(); err != nil {
				return err
			}

			if err := req.Update(); err != nil {
				return err
			}

			fmt.Printf("%q was successfully saved\n", teamID)
			return nil
		},
		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("team identifier must be set as an argument")
			}
			return nil
		},
	}
}

// GetDeleteTeamMemberCommand returns the delete team member command
func GetDeleteTeamMemberCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Removes a member from the given team",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "The username of the user you wish to remove from the team",
				Required: true,
			},
		},

		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")
			if team == "" {
				return errTeamParameterMissing
			}

			exists, err := ResourceExists(config, "team", team)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("%q does not exist", team)
			}

			err = NewRequest().
				WithConfig(config).
				WithContext(ctx).
				WithEndpoint("/teams/{team}/members/{user}").
				PathParameter("team", true).
				PathParameter("user", true).
				Delete()
			if err != nil {
				return err
			}
			fmt.Printf("%q has been successfully removed from team %q\n",
				ctx.String("user"), ctx.String("team"))

			return nil
		},
	}
}
