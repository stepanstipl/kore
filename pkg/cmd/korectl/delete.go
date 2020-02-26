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
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
)

func GetDeleteCommand(config *Config) cli.Command {
	return cli.Command{
		Name:      "delete",
		Aliases:   []string{"rm", "del"},
		Usage:     "Used to delete one or more resources from the kore",
		ArgsUsage: "-f <file> | <kind> <name>",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "file,f",
				Usage: "The path to the file containing the resources definitions `PATH`",
			},
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to filter the results by team `TEAM`",
			},
		},
		Action: func(ctx *cli.Context) error {
			for _, file := range ctx.StringSlice("file") {
				// @step: read in the content of the file
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				documents, err := ParseDocument(bytes.NewReader(content), ctx.String("team"))
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
				if ctx.NArg() != 2 {
					return errors.New("you need to specify a resource type and ")
				}
				req := NewRequest().
					WithConfig(config).
					WithContext(ctx).
					PathParameter("resource", true).
					PathParameter("name", false)

				endpoint := "/teams/{team}/{resource}/{name}"

				if IsGlobalResource(ctx.Args().First()) {
					endpoint = "/{resource}/{name}"
				} else {
					req.PathParameter("team", true)
				}
				req.WithEndpoint(endpoint).
					WithInject("resource", ctx.Args().First()).
					WithInject("name", ctx.Args().Tail()[0])

				if err := req.Delete(); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
