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
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	"github.com/urfave/cli/v2"
)

var (
	createNamespaceLongDescription = `
Provides the ability to create a namespace on a provisioned cluster. In order to
retrieve the clusters you have available you can run:

$ korectl get clusters -t <team> 

Examples:
# Create a namespace on cluster 'dev'
$ korectl create namespace -c cluster -t <team>

# Deleting a namespace on the cluster
$ korectl delete namespaceclaim

You can list the namespace you have already provisioned via

$ korectl get namespaceclaims -t <team>
`
)

// GetCreateNamespaceCommand creates and returns the create namespace command
func GetCreateNamespaceCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "namespace",
		Description: createNamespaceLongDescription,
		Usage:       "Create a namespace on the cluster",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cluster,c",
				Usage: "the name of the cluster you want the namespace to reside `NAME`",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "generate the cluster specification but does not apply `BOOL`",
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the namespace should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			cluster := ctx.String("cluster")
			team := ctx.String("team")
			dry := ctx.Bool("dry-run")

			// @step: evaluate the options
			if team == "" {
				return errTeamParameterMissing
			}
			if cluster == "" {
				return fmt.Errorf("you must specify a cluster: $ korectl get clusters -t %s", team)
			}

			// @step: check the kubernetes cluster exists
			if found, err := TeamResourceExists(config, team, "cluster", cluster); err != nil {
				return err
			} else if !found {
				return fmt.Errorf("cluster: %s does not exist", cluster)
			}

			owner := corev1.Ownership{
				Group:     clustersv1.GroupVersion.Group,
				Version:   clustersv1.GroupVersion.Version,
				Kind:      "Kubernetes",
				Namespace: team,
				Name:      cluster,
			}

			if err := CreateClusterNamespace(config, owner, team, name, dry); err != nil {
				return fmt.Errorf("trying to provision namespace on cluster: %s", err)
			}
			fmt.Println("Namespace provisioning on the cluster, you can check via: $ korectl get namespaceclaims -t", team)

			return nil
		},
	}
}
