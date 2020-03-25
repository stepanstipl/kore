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
	"fmt"

	"github.com/appvia/kore/pkg/utils"

	"github.com/urfave/cli/v2"
)

// GetApplyCommand returns the resource apply command
func GetApplyCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:  "apply",
		Usage: "Used to apply one of more resources to the API",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "path to the file containing resource definition/s (use '-' for stdin) `PATH`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			for _, file := range ctx.StringSlice("file") {
				// @step: read in the content of the file
				content, err := utils.ReadFileOrStdin(file)
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
