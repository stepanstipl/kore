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
	"net/http"

	"github.com/urfave/cli"
)

var getLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Examples:
  List users:
  $ korectl get users

  Get information about a specific user:
  $ korectl get user admin -o yaml
`

func GetGetCommand(config *Config) cli.Command {
	return cli.Command{
		Name:        "get",
		Usage:       "Retrieves one or more resources from the api",
		Description: formatLongDescription(getLongDescription),
		ArgsUsage:   "[TYPE] [NAME]",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to filter the results by team",
			},
		}, DefaultOptions...),

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return errors.New("you need to specify a resource type")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			req, resourceConfig, err := NewRequestForResource(config, ctx)
			if err != nil {
				return err
			}

			req.Render(resourceConfig.Columns...)

			if err := req.Get(); err != nil {
				if reqErr, ok := err.(*RequestError); ok {
					if reqErr.statusCode == http.StatusNotFound {
						if ctx.NArg() == 1 {
							return fmt.Errorf("%q is not a valid resource type", ctx.Args().Get(0))
						} else {
							return fmt.Errorf("%q does not exist", ctx.Args().Get(1))
						}
					}
				}
				return err
			}
			fmt.Println("")

			return nil
		},
	}
}
