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
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"
)

func GetApplyCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "apply",
		Usage: "Used to apply one of more resources to the API",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:     "file,f",
				Usage:    "The path to the file containing the resources definitions `PATH`",
				Required: true,
			},
			cli.StringFlag{
				Name:  "team,t",
				Usage: "Used to filter the results by team `TEAM`",
			},
		},
		Action: func(ctx *cli.Context) error {
			for _, file := range ctx.StringSlice("file") {
				// @step: get the options
				team := GetGlobalTeamFlag(ctx)

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
						Update()
					if err != nil {
						fmt.Printf("%s/%s failed with error: %s\n", gvk.Group, x.Endpoint, err)

						return err
					}

					fmt.Printf("%s/%s configured\n", gvk.Group, x.Endpoint)
				}
			}

			return nil
		},
	}
}
