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
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

func GetGetCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Used to retrieve a resource from the api",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to filter the results by team `TEAM`",
			},
		}, DefaultOptions...),

		Subcommands: []cli.Command{
			GetGetTeamCommand(config),
		},

		Action: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return errors.New("you need to specify a resource type")
			}
			req := NewRequest().
				WithConfig(config).
				WithContext(ctx).
				PathParameter("resource", true).
				PathParameter("name", false)

			endpoint := "/teams/{team}/{resource}/{name}"

			if IsGlobalResource(ctx.Args().First()) {
				endpoint = "/{resource}/{name}"
			} else if IsGlobalResourceOptional(ctx.Args().First()) {
				if !ctx.IsSet("team") {
					endpoint = "/{resource}/{name}"
				} else {
					req.PathParameter("team", true)
				}
			} else {
				req.PathParameter("team", true)
			}

			req.WithEndpoint(endpoint).
				WithInject("resource", ctx.Args().First()).
				Render(
					Column("Name", ".metadata.name"),
				)

			if ctx.NArg() == 2 {
				req.WithInject("name", ctx.Args().Get(1))
			}
			if err := req.Get(); err != nil {
				return err
			}
			fmt.Println("")

			return nil
		},
	}
}
