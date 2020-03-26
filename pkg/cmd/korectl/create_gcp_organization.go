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

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/kore"

	"github.com/urfave/cli/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	gcpOrganzationgLongDescription = `
When using GCP as a cloud provider, kore makes use of a administrative
admin project as a means to manage team project and can assigned dedicated
projects to your teams.

# Check what organization are managed (as an admin)
$ korectl get organizations

Administrators can also assign gcp organization to you via allocations

$ korectl get allocations -t <team>

# Create a GCP organization (please refer to documentation for permissions).

We first create a secret holding the service account json

$ korectl create secret gcp \
  --type=gcp
	--from-file=key=<path>

$ korectl create google organization <name> \
	--parent-type=organization \
	--parent-id=11111111111 \
	--billing-account-id=018ACC-AD48A4-42ADD2 \
	--credentials=gcp
	-t kore-admin

Once created your check the permissions where correctly via

$ korectl get organizations <name> -t yaml -t kore-admin

We can then allocate the credentials out to teams via: korectl get allocations --help
`
)

// GetCreateGCPOrganization provides the gcp create commands
func GetCreateGCPOrganization(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "organization",
		Aliases:     []string{"org"},
		Usage:       "is used to provision an admistrative project to provide team projects",
		Description: formatLongDescription(gcpOrganzationgLongDescription),
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "parent-type",
				Usage: "Is the gcp parent of the project i.e project, org etc `NAME`",
				Value: "organization",
			},
			&cli.StringFlag{
				Name:     "parent-id",
				Usage:    "Is the parent id of the above, i.e the project id or the organizational id `ID`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "billing-account-id",
				Usage:    "is the billing account id which allow teams projects should get billing `ID`",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "service-account-name",
				Usage: "the default service account name to use in the admin project `NAME`",
				Value: "kore",
			},
			/*
				&cli.StringFlag{
					Name:  "oauth-token",
					Usage: "The name of the kore secret which holds the oauth token `NAME`",
				},
			*/
			&cli.StringFlag{
				Name:  "credentials",
				Usage: "The name of the kore secret holding the service account key `NAME`",
			},
			&cli.StringFlag{
				Name:  "credentials-team",
				Usage: "The name of the team whom owns the credentials `NAME`",
				Value: kore.HubAdminTeam,
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the organization should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			team := ctx.String("team")
			name := ctx.Args().First()
			kind := "organization"
			cname := ctx.String("credentials")
			cnamespace := ctx.String("credentials-team")
			nowait := ctx.Bool("no-wait")

			found, err := TeamResourceExists(config, team, kind, name)
			if err != nil {
				return err
			}
			if found {
				return fmt.Errorf("%q already exists, please edit instead", name)
			}

			o := &gcp.Organization{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: team,
				},
				Spec: gcp.OrganizationSpec{
					BillingAccount: ctx.String("billing-account-id"),
					ParentID:       ctx.String("parent-id"),
					ParentType:     ctx.String("parent-type"),
					ServiceAccount: ctx.String("service-account-name"),
					CredentialsRef: &v1.SecretReference{
						Name:      cname,
						Namespace: cnamespace,
					},
				},
			}

			if err := CreateTeamResource(config, team, kind, name, o); err != nil {
				return err
			}

			return WaitForResourceCheck(context.Background(), config, team, kind, name, nowait)
		},
	}
}
