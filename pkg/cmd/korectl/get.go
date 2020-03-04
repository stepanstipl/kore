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

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return errors.New("you need to specify a resource type")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			// @step: setup the printer for the resource type
			printer := resourcePrinters.Get(ctx.Args().First())

			req := NewRequest().
				WithConfig(config).
				WithContext(ctx).
				PathParameter("resource", true).
				PathParameter("name", false).
				Render(printer.Columns...).
				WithInject("resource", printer.APIResourceName)

			endpoint := "/teams/{team}/{resource}/{name}"

			// @check if the resource is a global resource i.e plans, teams, users etc
			switch {
			case IsGlobalResource(printer.APIResourceName):
				endpoint = "/{resource}/{name}"

			case IsGlobalResourceOptional(printer.APIResourceName):
				switch ctx.IsSet("team") {
				case true:
					req.WithInject("team", GlobalStringFlag(ctx, "team"))
					req.PathParameter("team", true)
				default:
					endpoint = "/{resource}/{name}"
				}

			default:
				req.PathParameter("team", true)
			}

			// @step: check if we are getting a specific resource by name
			if ctx.NArg() == 2 {
				req.WithInject("name", ctx.Args().Get(1))
			}

			// @step: retrieve the resource or resources from the api
			if err := req.WithEndpoint(endpoint).Get(); err != nil {
				return err
			}
			fmt.Println("")

			return nil
		},
	}
}
