/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import "github.com/urfave/cli"

// GetUsersCommand returns the users command
func GetUsersCommands(config *Config) cli.Command {
	return cli.Command{
		Name:    "users",
		Aliases: []string{"us"},
		Usage:   "Used to get, list, update and delete users in the hub",

		Subcommands: []cli.Command{
			{
				Name:  "get",
				Usage: "Used to retrieve one of more users from the hub",
				Flags: append([]cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the user to retrieve `NAME`",
					},
				}, DefaultOptions...),
				Action: func(ctx *cli.Context) error {
					return NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/users").
						PathParameter("name", false).
						Render(
							Column("Username", ".spec.username"),
							Column("Email", ".spec.email"),
						).
						Get()
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"rm"},
				Usage:   "Used to delete a user from the hub",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "name,n",
						Usage:    "The name of the user to remove `NAME`",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					err := NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/users/{name}").
						PathParameter("name", true).
						Delete()
					if err != nil {
						return err
					}
					return NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/users").
						Render(
							Column("Username", ".spec.username"),
							Column("Email", ".spec.email"),
						).
						Get()
				},
			},
		},
	}
}
