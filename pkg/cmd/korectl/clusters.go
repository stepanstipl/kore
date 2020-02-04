/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/utils"

	"github.com/urfave/cli"
)

func GetClustersCommand(config *Config) cli.Command {
	return cli.Command{
		Name:    "clusters",
		Aliases: []string{"cls"},
		Usage:   "Used to manage and interact with clusters provisioned by the kore",
		Subcommands: []cli.Command{
			{
				Name:  "auth",
				Usage: "Used to retrieve the API endpoints of the clusters and provision your kubeconfig",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
					cli.StringFlag{
						Name:  "team,t",
						Usage: "Used to filter the results by team `TEAM`",
					},
				},
				Action: func(ctx *cli.Context) error {
					resp := NewRequest()
					err := resp.
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/teams/{team}/clusters").
						PathParameter("team", true).
						Get()
					if err != nil {
						return err
					}
					// else we need to provision the kubeconfig
					clusters := &clustersv1.KubernetesList{}
					if err := utils.DecodeToJSON(resp.Body(), clusters); err != nil {
						return err
					}
					if len(clusters.Items) <= 0 {
						fmt.Println("no clusters found in this teams namespace")

						return nil
					}

					kubeconfig, err := GetKubeConfig()
					if err != nil {
						return err
					}

					if err := PopulateKubeconfig(clusters, kubeconfig, config); err != nil {
						return err
					}
					fmt.Println("Successfull updated you kubeconfig with credentuals")

					return nil
				},
			},
			{
				Name:  "get",
				Usage: "Used to retrieve one of all cluster from the kore",
				Flags: append([]cli.Flag{
					cli.StringFlag{
						Name:  "name,n",
						Usage: "The name of the integration to retrieve `NAME`",
					},
					cli.StringFlag{
						Name:  "team,t",
						Usage: "Used to filter the results by team `TEAM`",
					},
				}, DefaultOptions...),
				Action: func(ctx *cli.Context) error {
					return NewRequest().
						WithConfig(config).
						WithContext(ctx).
						WithEndpoint("/teams/{team}/clusters").
						PathParameter("team", true).
						PathParameter("name", false).
						Render(
							Column("Name", ".metadata.name"),
							Column("Domain", ".spec.domain"),
							Column("Endpoint", ".status.endpoint"),
							Column("Status", ".status.status"),
						).
						Get()
				},
			},
		},
	}
}
