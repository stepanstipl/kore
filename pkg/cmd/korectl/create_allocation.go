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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"

	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createAllocationLongDescription = `
Allocations are used to share resources (secrets) functionality (cloud, clusters)
to one or more teams in kore. This allocation are then used as references to the
containing resource i.e. a namespace refers to a cluster, a cloud provider refers
to credentials and so forth.

# Create an allocation for a gcp organization to all teams
$ korectl create allocation <name> -t <team> \
	--name=gcp
	--group=gcp.compute.appvia.io \
	--version=v1alpha1 \
	--kind=Organization \
	--to=*

# Retrieve a list all allocations
$ korectl get allocations -t <team>
`
)

// GetCreateAllocation provides the create allocation command
func GetCreateAllocation(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "allocation",
		Aliases:     []string{"allocations"},
		Usage:       "allows you to create an allocation for teams to consume",
		Description: formatLongDescription(createAllocationLongDescription),
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "description",
				Usage:    "provide a description defining what you are allocating `DESCRIPTION`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "resource",
				Usage:    "is the name of 'resource' you are allocation (not the name of the allocation) `NAME`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "group",
				Usage:    "is the API group the resource being allocated resides `NAME`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "version",
				Usage:    "is API version the resource being allocated resides in e.g v1 `VERSION`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "kind",
				Usage:    "the resource kind (kubernetes type) which is being allocated  `NAME`",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "to",
				Usage: "the teams the resource is being allocated to (can use more than once) `TEAM`",
				Value: cli.NewStringSlice("*"),
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the allocation should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")
			name := ctx.Args().First()

			description := ctx.String("description")
			group := ctx.String("group")
			kind := ctx.String("kind")
			resource := ctx.String("resource")
			teams := ctx.StringSlice("to")
			version := ctx.String("version")
			nowait := ctx.Bool("no-wait")

			found, err := TeamResourceExists(config, team, "allocation", name)
			if err != nil {
				return err
			}
			if found {
				return fmt.Errorf("%q already exists, please edit instead", name)
			}

			o := &configv1.Allocation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: kore.HubAdminTeam,
				},
				Spec: configv1.AllocationSpec{
					Name:    name,
					Summary: description,
					Resource: corev1.Ownership{
						Group:     group,
						Kind:      kind,
						Name:      resource,
						Namespace: kore.HubAdminTeam,
						Version:   version,
					},
					Teams: teams,
				},
			}

			if err := CreateTeamResource(config, team, "allocation", name, o); err != nil {
				return err
			}

			return WaitForResourceCheck(context.Background(), config, team, kind, name, nowait)
		},
	}
}
