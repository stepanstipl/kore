/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"fmt"
	"github.com/urfave/cli"
)

// GetTeamsCommand returns the teams command
func GetTeamsCommands(config *Config) cli.Command {
	return cli.Command{
		Name:    "teams",
		Aliases: []string{"tm"},
		Usage:   "Used to interact, get, list and update teams in the kore",
		Subcommands: []cli.Command{
			{
				Name:  "get",
				Usage: "Used to retrieve the details of a team in the kore",
				Flags: append([]cli.Flag{
					cli.StringFlag{
						Name:     "name,n",
						Usage:    "The name of the team to retrieve (assumes all if not defined)",
						Required: false,
					},
				}, DefaultOptions...),
				Action: func(c *cli.Context) error {
					return NewRequest().
						WithConfig(config).
						WithContext(c).
						WithEndpoint("/teams/{name}").
						PathParameter("name", false).
						Render(
							Column("Name", ".metadata.name"),
							Column("Description", ".spec.description"),
						).
						Get()
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"rm"},
				Usage:   "Used to delete a team from the kore",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the team to remove",
					},
				},
				Action: func(ctx *cli.Context) error {
					return NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/teams/{name}").
						PathParameter("name", true).
						Delete()
				},
			},
			{
				Name:    "members",
				Aliases: []string{"mb"},
				Usage:   "Used to get, list, add and remove users to the team",
				Subcommands: []cli.Command{
					{
						Name:  "get",
						Usage: "Used to list all the users in the team",
						Flags: append([]cli.Flag{
							cli.StringFlag{
								Name:  "team,t",
								Usage: "The name of the team you wish to list the users",
							},
						}, DefaultOptions...),
						Action: func(ctx *cli.Context) error {
							return NewRequest().
								WithConfig(config).
								WithContext(ctx).
								WithEndpoint("/teams/{team}/members").
								PathParameter("team", true).
								Render(
									Column("Username", "."),
								).
								Get()
						},
					},
					{
						Name:  "add",
						Usage: "Used to add a kore member into the team",
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:     "team,t",
								Usage:    "The name of the team you wish to list the users",
								Required: true,
							},
							cli.StringFlag{
								Name:     "user,u",
								Usage:    "The name of the user you wish to add to the team",
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
							fmt.Printf("[status] user %s has been added to team: %s\n", ctx.String("user"), ctx.String("team"))

							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "Used to remove a member from a team in th kore",
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:     "team,t",
								Usage:    "The name of the team you wish to remove the user from `TEAM`",
								Required: true,
							},
							cli.StringFlag{
								Name:     "user,u",
								Usage:    "The name of the user you wish to remove from the team `USERNAME`",
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
							fmt.Printf("[status] user %s has been remove to team: %s\n",
								ctx.String("user"), ctx.String("team"))

							return nil
						},
					},
				},
			},
		},
	}
}
