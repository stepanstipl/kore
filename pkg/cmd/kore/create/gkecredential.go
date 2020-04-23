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

package create

import (
	"fmt"
	"io/ioutil"

	confv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateGKECredentialsOptions is used to provision a team
type CreateGKECredentialsOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the name of the credential
	Name string
	// DryRun indicates we only dryrun the resources
	DryRun bool
	// Description is a description of the credential
	Description string
	// ProjectName is the name of the GCP project
	ProjectName string
	// ServiceAccountJSON is a reference to a file containing the JSON service account details.
	ServiceAccountJSON string
	// AllocateToTeams allows the credential to be allocated to the specified list of teams.
	AllocateToTeams []string
	// AllocateToAll controls if a default allocation should be set for this to allocate to all teams.
	AllocateToAll bool
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
}

// NewCmdGKECredentials returns the create GCP project credentials command
func NewCmdGKECredentials(factory cmdutil.Factory) *cobra.Command {
	o := &CreateGKECredentialsOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "gkecredentials",
		Aliases: []string{"gkecredential"},
		Short:   "Creates a set of GKE project-level credentials in Kore",
		Example: "kore create gkecredentials <name> -d <description> -p <gcp project> --cred-file ./service-account.json",
		PreRunE: cmdutil.RequireName,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Description, "description", "d", "", "the description of the credential")
	flags.StringVarP(&o.ProjectName, "project", "p", "", "the GCP project for these credentials")
	flags.StringVarP(&o.ServiceAccountJSON, "cred-file", "c", "", "the service account JSON file containing the credentials to import")
	flags.StringArrayVarP(&o.AllocateToTeams, "allocate", "a", []string{}, "list of teams to allocate to, e.g. team1,team2")
	flags.BoolVar(&o.AllocateToAll, "all-teams", false, "make these credentials available to all teams in kore (if not set, you must create an allocation for these credentials for them to be usable)")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutil.MustMarkFlagRequired(command, "project")
	cmdutil.MustMarkFlagRequired(command, "cred-file")

	return command
}

// Run is responsible for creating the credentials
func (o CreateGKECredentialsOptions) Run() error {
	json, err := ioutil.ReadFile(o.ServiceAccountJSON)
	if err != nil {
		return fmt.Errorf("trying reading service account from %v", o.ServiceAccountJSON)
	}

	secret := &confv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: confv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: confv1.SecretSpec{
			Type:        "gke-credentials",
			Description: o.Description,
			Data: map[string]string{
				"service_account_key": string(json),
			},
		},
	}
	secret.Encode()

	o.Println("Storing credentials secret in Kore")
	err = o.WaitForCreation(
		o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("secret")).
			Name(o.Name).
			Payload(secret).
			Result(&confv1.Secret{}),
		o.NoWait,
	)
	if err != nil {
		return fmt.Errorf("trying to create secret: %s", err)
	}

	cred := &gke.GKECredentials{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GKECredentials",
			APIVersion: gke.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: gke.GKECredentialsSpec{
			Project: o.ProjectName,
			CredentialsRef: &v1.SecretReference{
				Name:      o.Name,
				Namespace: kore.HubAdminTeam,
			},
		},
	}

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Format(render.FormatYAML).
			Resource(render.FromStruct(cred)).
			Do()
	}

	err = o.WaitForCreation(
		o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("gkecredential")).
			Name(o.Name).
			Payload(cred).
			Result(&gke.GKECredentials{}),
		o.NoWait,
	)
	if err != nil {
		return fmt.Errorf("trying to create credential: %s", err)
	}

	if !o.AllocateToAll && len(o.AllocateToTeams) == 0 {
		return nil
	}

	// Create allocation
	teams := o.AllocateToTeams
	if o.AllocateToAll {
		teams = []string{"*"}
	}
	alloc := &confv1.Allocation{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Allocation",
			APIVersion: confv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: confv1.AllocationSpec{
			Name:    o.Name,
			Summary: o.Description,
			Teams:   teams,
			Resource: corev1.Ownership{
				Group:     gke.GroupVersion.Group,
				Version:   gke.GroupVersion.Version,
				Kind:      "GKECredentials",
				Name:      o.Name,
				Namespace: kore.HubAdminTeam,
			},
		},
	}

	o.Println("Storing credential allocation in Kore")

	return o.WaitForCreation(
		o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("allocation")).
			Name(o.Name).
			Payload(alloc).
			Result(&confv1.Allocation{}),
		o.NoWait,
	)
}
