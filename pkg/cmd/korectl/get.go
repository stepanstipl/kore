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
	"errors"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

var getLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Examples:
  List users:
  $ korectl get users

  Get information about a specific user:
  $ korectl get user admin -o yaml
`

func GetGetCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "Retrieves one or more resources from the api",
		Description: formatLongDescription(getLongDescription),
		ArgsUsage:   "[TYPE] [NAME]",
		Flags: append([]cli.Flag{
			&cli.StringFlag{
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
