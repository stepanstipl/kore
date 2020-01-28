/**
 * Copyright (C) 2020 Rohith Jayawardene <gambol99@gmail.com>
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

package hubctl

import "github.com/urfave/cli"

func GetIntegrationCommands() cli.Command {
	return cli.Command{
		Name:    "integrations",
		Aliases: []string{"int"},
		Usage:   "Provides the ability to interact and managed the integrations within the hub",
		Subcommands: []cli.Command{
			{
				Name:  "get",
				Usage: "Retrieves one of all the integrations configured witin the hub",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
				},
			},
			{
				Name:  "apply",
				Usage: "Used to apply an integration to the hub",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
					cli.StringFlag{
						Name:  "file,f",
						Usage: "The path to the file containing the integration definition `PATH`",
					},
				},
			},
			{
				Name:  "edit",
				Usage: "Permits you to edit the intgration in the hub",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "name,n",
						Usage:    "The name of the integration to edit `NAME`",
						Required: true,
					},
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "Used to remove an integeration from the hub",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to edit `NAME`",
					},
					cli.StringFlag{
						Name:  "file,f",
						Usage: "The path to the file containing the integration definition `PATH`",
					},
				},
			},
		},
	}
}
