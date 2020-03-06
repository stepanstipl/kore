/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"fmt"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetCreateTeamCommand(config *Config) cli.Command {
	return cli.Command{
		Name:      "team",
		Aliases:   []string{"teams"},
		Usage:     "Creates a team",
		ArgsUsage: "TEAM",
		Flags: []cli.Flag{
			cli.StringFlag{
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

func GetEditTeamCommand(config *Config) cli.Command {
	return cli.Command{
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

func GetCreateTeamMemberCommand(config *Config) cli.Command {
	return cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Creates a new team member",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "team,t",
				Usage:    "The name of the team you wish to add the user to",
				Required: true,
			},
			cli.StringFlag{
				Name:     "user,u",
				Usage:    "The username of the user you wish to add to the team",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			err := NewRequest().
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

func GetDeleteTeamMemberCommand(config *Config) cli.Command {
	return cli.Command{
		Name:    "member",
		Aliases: []string{"members"},
		Usage:   "Removes a member from the given team",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "team,t",
				Usage:    "The name of the team you wish to remove the user from",
				Required: true,
			},
			cli.StringFlag{
				Name:     "user,u",
				Usage:    "The username of the user you wish to remove from the team",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			err := NewRequest().
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
