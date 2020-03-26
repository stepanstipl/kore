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

// GetGetCommand returns the auto-generated resources
func GetGetCommand(config *Config) *cli.Command {
	// @note: the moment the user need to 'know' the name of all resource
	// types for them to retrieve from the API. Ideally this would be retrieved
	// from the API (later to come) but for now to complete the story we need
	// a means to these to autocomplete on the CLI to make it apparent.

	var commands []*cli.Command

	for k, v := range resourceConfigs {
		usage := "retrieve the " + k + " resource"
		if v.IsGlobal {
			usage = "retrieve the global resource " + k
		}

		command := &cli.Command{
			Name:    k,
			Aliases: []string{v.Name},
			Usage:   usage,

			Action: func(ctx *cli.Context) error {
				team := ctx.String("team")
				name := ctx.Args().First()
				resource := getResourceConfig(ctx.Command.Name)

				// @step: determine the resource request type
				request, _, err := NewRequestForResource(config, team, resource.Name, name)
				if err != nil {
					return err
				}
				request = request.Render(resource.Columns...).WithContext(ctx)

				// @step: make the request to the api
				if err := request.Get(); err != nil {
					if reqErr, ok := err.(*RequestError); ok {
						if reqErr.statusCode == http.StatusNotFound {
							if ctx.NArg() == 1 {
								return fmt.Errorf("%q is not a valid resource type", ctx.Args().Get(0))
							}

							return fmt.Errorf("%q does not exist", ctx.Args().Get(1))
						}
					}

					return err
				}
				fmt.Println("")

				return nil
			},
		}

		commands = append(commands, command)
	}

	return &cli.Command{
		Name:        "get",
		Usage:       "Retrieves one or more resources from the api",
		Description: formatLongDescription(getLongDescription),
		ArgsUsage:   "[TYPE] [NAME]",
		Subcommands: commands,
	}
}
