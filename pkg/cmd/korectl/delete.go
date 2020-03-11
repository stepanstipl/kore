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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli/v2"
)

var deleteLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Example to delete a user:
  $ korectl delete user joe@example.com

Example to delete multiple resources from a file:
  $ korectl delete --file resources.yaml
`

// GetDeleteCommand creates and returns the delete command
func GetDeleteCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Aliases:     []string{"rm", "del"},
		Usage:       "Deletes one or more resources",
		Description: formatLongDescription(deleteLongDescription),
		ArgsUsage:   "-f <file> | [TYPE] [NAME]",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "The path to the file containing the resources definitions `PATH`",
			},
		},
		Subcommands: []*cli.Command{
			GetDeleteTeamMemberCommand(config),
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.IsSet("file") && !ctx.Args().Present() {
				_ = cli.ShowSubcommandHelp(ctx)
				return fmt.Errorf("-f <file> or [TYPE] [NAME] is required")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")

			for _, file := range ctx.StringSlice("file") {
				// @step: read in the content of the file
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				documents, err := ParseDocument(bytes.NewReader(content), team)
				if err != nil {
					return err
				}
				for _, x := range documents {
					gvk := x.Object.GetObjectKind().GroupVersionKind()
					err := NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint(x.Endpoint).
						WithRuntimeObject(x.Object).
						Delete()
					if err != nil {
						fmt.Printf("%s/%s failed with error: %s\n", gvk.Group, x.Endpoint, err)

						return err
					}

					fmt.Printf("%s/%s deleted\n", gvk.Group, x.Endpoint)
				}
			}
			if len(ctx.StringSlice("file")) <= 0 {
				if ctx.NArg() < 2 {
					return errors.New("you need to specify a resource type and name")
				}

				req, _, err := NewRequestForResource(config, ctx)
				if err != nil {
					return err
				}

				exists, err := req.Exists()
				if err != nil {
					return err
				}

				if !exists {
					return fmt.Errorf("%q does not exist", ctx.Args().Get(1))
				}

				if err := req.Delete(); err != nil {
					return err
				}

				fmt.Printf("%q was successfully deleted\n", ctx.Args().Get(1))
			}

			return nil
		},
	}
}
