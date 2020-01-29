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

package korectl

import "github.com/urfave/cli"

func GetWhoamiCommand(config *Config) cli.Command {
	return cli.Command{
		Name:    "whoami",
		Aliases: []string{"who"},
		Usage:   "Used to retrieve details on your identity within the hub",

		Action: func(ctx *cli.Context) error {
			return NewRequest().
				WithConfig(config).
				WithContext(ctx).
				WithEndpoint("/whoami").
				Render(
					Column("Username", ".username"),
					Column("Email", ".email"),
					Column("Teams", ".teams"),
				).
				Get()
		},
	}
}
