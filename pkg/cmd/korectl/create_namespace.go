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
	"context"
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
		Description: formatLongDescription(createNamespaceLongDescription),
		Usage:       "Create a namespace on the cluster",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cluster",
				Aliases: []string{"c"},
				Usage:   "the name of the cluster you want the namespace to reside `NAME`",
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
			dry := ctx.Bool("dry-run")
			kind := "namespaceclaim"
			team := ctx.String("team")
			wait := ctx.Bool("wait")

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

			return WaitForResourceCheck(context.Background(), config, team, kind, name, wait)
		},
	}
}
