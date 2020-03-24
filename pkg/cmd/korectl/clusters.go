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
	"github.com/urfave/cli/v2"
)

func GetClustersCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "clusters",
		Aliases: []string{"cls"},
		Usage:   "Used to manage and interact with clusters provisioned by the kore",
		Subcommands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Used to retrieve one or all clusters from the kore",
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "The name of the integration to retrieve `NAME`",
					},
				}, DefaultOptions...),
				Action: func(ctx *cli.Context) error {
					team := ctx.String("team")

					return NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/teams/{team}/clusters").
						WithInject("team", team).
						PathParameter("team", true).
						PathParameter("name", false).
						Render(
							Column("Name", "metadata.name"),
							Column("Endpoint", "status.endpoint"),
							Column("Status", "status.status"),
						).
						Get()
				},
			},
		},
	}
}
