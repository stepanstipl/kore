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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	"github.com/urfave/cli/v2"
)

func GetClustersCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:    "clusters",
		Aliases: []string{"cls"},
		Usage:   "Used to manage and interact with clusters provisioned by the kore",
		Subcommands: []*cli.Command{
			{
				Name:  "auth",
				Usage: "Used to retrieve the API endpoints of the clusters and provision your kubeconfig",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
					&cli.StringFlag{
						Name:  "team,t",
						Usage: "Used to filter the results by team `TEAM`",
					},
				},
				Action: func(ctx *cli.Context) error {
					clusters := &clustersv1.KubernetesList{}
					team := ctx.String("team")

					if err := GetTeamResourceList(config, team, "clusters", clusters); err != nil {
						return err
					}

					if len(clusters.Items) <= 0 {
						fmt.Println("no clusters found in this team's namespace")

						return nil
					}

					kubeconfig, err := GetKubeConfig()
					if err != nil {
						return err
					}

					if err := PopulateKubeconfig(clusters, kubeconfig, config); err != nil {
						return err
					}
					fmt.Println("Successfully updated your kubeconfig with credentials")

					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Used to retrieve one or all clusters from the kore",
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
					&cli.StringFlag{
						Name:  "team,t",
						Usage: "Used to filter the results by team `TEAM`",
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
							Column("Name", ".metadata.name"),
							Column("Endpoint", ".status.endpoint"),
							Column("Status", ".status.status"),
						).
						Get()
				},
			},
		},
	}
}
