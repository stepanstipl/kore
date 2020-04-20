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

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateGKECredentialsOptions is used to provision a team
type CreateGKECredentialsOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the name of the credential
	Name string
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

	command.Flags().StringVarP(&o.Description, "description", "d", "", "the description of the credential")
	command.Flags().StringVarP(&o.ProjectName, "project", "p", "", "the GCP project for these credentials")
	command.Flags().StringVarP(&o.ServiceAccountJSON, "cred-file", "c", "", "the service account JSON file containing the credentials to import")
	command.Flags().StringArrayVarP(&o.AllocateToTeams, "allocate", "a", []string{}, "list of teams to allocate to, e.g. team1,team2")
	command.Flags().BoolVar(&o.AllocateToAll, "all-teams", false, "make these credentials available to all teams in kore (if not set, you must create an allocation for these credentials for them to be usable)")

	cmdutil.MustMarkFlagRequired(command, "project")
	cmdutil.MustMarkFlagRequired(command, "cred-file")

	return command
}

// Run is responsible for creating the credentials
func (o CreateGKECredentialsOptions) Run() error {
	found, err := o.Client().Team(kore.HubAdminTeam).Resource("gkecredential").Name(o.Name).Exists()
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("%q already exists, please edit instead", o.Name)
	}

	json, err := ioutil.ReadFile(o.ServiceAccountJSON)
	if err != nil {
		o.Println("Error reading service account from %v", o.ServiceAccountJSON)

		return err
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
			Account: string(json),
			Project: o.ProjectName,
		},
	}

	o.Println("Storing credentials in Kore")
	err = o.WaitForCreation(
		o.Client().
			Team(kore.HubAdminTeam).
			Resource("gkecredential").
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
		o.Client().
			Team(kore.HubAdminTeam).
			Resource("allocation").
			Name(o.Name).
			Payload(alloc).
			Result(&confv1.Allocation{}),
		o.NoWait,
	)
}
