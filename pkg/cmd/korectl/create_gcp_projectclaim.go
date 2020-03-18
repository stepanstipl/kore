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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"

	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	gcpProjectLongDescription = `
When using GCP as a cloud provider, teams are managed inside their own
GCP projects, isolating costs, risk, potential impact of upgrades
and maintenance. Administrators create and share out GCP Organizations to
one or more teams which the teams can then manage for themselves while
staying within the bounds on the orgs policy.

# Retrieve a list of organizations allocated to me
$ korectl get allocations -t <team>

# Request a project for my team
$ korectl create gcp project <name> --organization <allocation_name>

You can check the status of the project via
$ korectl get projectclaims -t devs
`
)

// GetCreateGCPProject provides the action to create team projects
func GetCreateGCPProject(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "project",
		Aliases:     []string{"projects"},
		Usage:       "provisions teams a gcp project to contain the resources",
		Description: formatLongDescription(gcpProjectLongDescription),
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "organization",
				Usage:    "the name of the gcp organization which the project should reside `NAME`",
				Aliases:  []string{"org"},
				Required: true,
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the project should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")
			name := ctx.Args().First()
			kind := "projectclaim"
			org := ctx.String("organization")

			found, err := TeamResourceExists(config, team, kind, name)
			if err != nil {
				return err
			}
			if found {
				return fmt.Errorf("%q already exists, please edit instead", name)
			}

			// @step: we retrieve the allocation
			found, err = TeamResourceExists(config, team, "allocation", org)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("%q allocation does not exist", org)
			}

			a, err := GetTeamAllocation(config, team, org)
			if err != nil {
				return err
			}

			// @step: check this is for a organization
			if a.Spec.Resource.Kind != "Organization" {
				return fmt.Errorf("%q allocation is not a gcp organization", org)
			}

			o := &gcp.ProjectClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: team,
				},
				Spec: gcp.ProjectClaimSpec{
					Organization: corev1.Ownership{
						Group:     a.Spec.Resource.Group,
						Kind:      a.Spec.Resource.Kind,
						Name:      a.Spec.Resource.Name,
						Namespace: a.Spec.Resource.Namespace,
						Version:   a.Spec.Resource.Version,
					},
				},
			}

			if err := CreateTeamResource(config, team, kind, name, o); err != nil {
				return err
			}
			fmt.Printf("%q has been successfully created\n", name)

			return nil
		},
	}
}
