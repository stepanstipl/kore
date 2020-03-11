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

	"github.com/manifoldco/promptui"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetCreateTeamCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:      "team",
		Aliases:   []string{"teams"},
		Usage:     "Creates a team",
		ArgsUsage: "TEAM",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "description",
				Usage:    "The description of the team",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			teamID := ctx.Args().First()

			exists, err := NewRequest().
				WithConfig(config).
				WithContext(ctx).
				PathParameter("id", true).
				WithInject("id", teamID).
				WithEndpoint("/teams/{id}").
				Exists()
			if err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("%q already exists", teamID)
			}

			team := &orgv1.Team{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Team",
					APIVersion: orgv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      teamID,
					Namespace: "",
					Labels:    nil,
				},
				Spec: orgv1.TeamSpec{
					Summary:     teamID,
					Description: ctx.String("description"),
				},
			}

			err = NewRequest().
				WithConfig(config).
				WithContext(ctx).
				PathParameter("id", true).
				WithInject("id", teamID).
				WithEndpoint("/teams/{id}").
				WithRuntimeObject(team).
				Update()
			if err != nil {
				return err
			}

			fmt.Printf("%q team was successfully created\n", teamID)
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

func GetCreateTeamMemberCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Creates a new team member",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user,u",
				Usage:    "The username of the user you wish to add to the team",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "invite,i",
				Usage:    "If the user doesn't exist and the invite flag is set, the invite url will be automatically generated.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")
			if team == "" {
				return errTeamParameterMissing
			}

			teamExists, err := ResourceExists(config, "team", team)
			if err != nil {
				return err
			}
			if !teamExists {
				return fmt.Errorf("team %q does not exist", team)
			}

			userExists, err := ResourceExists(config, "user", ctx.String("user"))
			if err != nil {
				return err
			}

			if !userExists {
				if !ctx.Bool("invite") {
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

				var inviteUrl string
				err := NewRequest().
					WithConfig(config).
					WithContext(ctx).
					WithEndpoint("/teams/{team}/invites/generate/{user}").
					PathParameter("team", true).
					WithInject("team", team).
					PathParameter("user", true).
					WithRuntimeObject(&inviteUrl).
					Get()
				if err != nil {
					return err
				}
				fmt.Printf("Invite URL: %s\n", inviteUrl)
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

func GetDeleteTeamMemberCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Removes a member from the given team",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "team,t",
				Usage: "The name of the team you wish to remove the user from",
			},
			&cli.StringFlag{
				Name:     "user,u",
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
