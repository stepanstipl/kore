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
	"strings"

	confv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateEKSCredentialsOptions is used to provision a team
type CreateEKSCredentialsOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the name of the credential
	Name string
	// Description is a description of the credential
	Description string
	// AccountID is the AWS numerical account ID
	AccountID string
	// AccessKeyID is the ID for the AWS access key
	AccessKeyID string
	// SecretAccessKey is the secret AWS access key
	SecretAccessKey string
	// AllocateToTeams allows the credential to be allocated to the specified list of teams.
	AllocateToTeams []string
	// AllocateToAll controls if a default allocation should be set for this to allocate to all teams.
	AllocateToAll bool
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
}

// NewCmdEKSCredentials returns the create EKS credentials command
func NewCmdEKSCredentials(factory cmdutil.Factory) *cobra.Command {
	o := &CreateEKSCredentialsOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "ekscredentials",
		Aliases: []string{"ekscredential"},
		Short:   "Creates a set of EKS / AWS credentials in Kore",
		Example: "kore create ekscredentials <name> -d <description> -i <aws account id> -k <aws access key id> -s <aws secret key> --all-teams",
		PreRunE: cmdutil.RequireName,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	command.Flags().StringVarP(&o.Description, "description", "d", "", "the description of the credential")
	command.Flags().StringVarP(&o.AccountID, "account-id", "i", "", "the AWS numerical account ID")
	command.Flags().StringVarP(&o.AccessKeyID, "key-id", "k", "", "the AWS access key ID")
	command.Flags().StringVarP(&o.SecretAccessKey, "secret-key", "s", "", "the AWS secret access key - will prompt if not provided")

	command.Flags().StringArrayVarP(&o.AllocateToTeams, "allocate", "a", []string{}, "list of teams to allocate to, e.g. team1,team2")
	command.Flags().BoolVar(&o.AllocateToAll, "all-teams", false, "make these credentials available to all teams in kore (if not set, you must create an allocation for these credentials for them to be usable)")

	cmdutil.MustMarkFlagRequired(command, "account-id")
	cmdutil.MustMarkFlagRequired(command, "key-id")

	return command
}

// Run is responsible for creating the credentials
func (o CreateEKSCredentialsOptions) Run() error {
	found, err := o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("ekscredential")).Name(o.Name).Exists()
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("%q already exists, please edit instead", o.Name)
	}

	if o.SecretAccessKey == "" {
		runner := promptui.Prompt{
			Label:   "Secret Access Key",
			Default: "",
			Validate: func(in string) error {
				if len(in) == 0 {
					return fmt.Errorf("Secret access key must be set")
				}
				return nil
			},
		}
		key, err := runner.Run()
		if err != nil {
			return err
		}
		o.SecretAccessKey = strings.TrimSpace(key)
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
			Type:        "aws-credential",
			Description: o.Description,
			Data: map[string]string{
				"access_id":     o.AccessKeyID,
				"access_secret": o.SecretAccessKey,
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
		return fmt.Errorf("Error while creating credential secret: %s", err)
	}

	cred := &eks.EKSCredentials{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EKSCredentials",
			APIVersion: eks.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubAdminTeam,
		},
		Spec: eks.EKSCredentialsSpec{
			AccountID: o.AccountID,
			CredentialsRef: &v1.SecretReference{
				Name:      o.Name,
				Namespace: kore.HubAdminTeam,
			},
		},
	}

	o.Println("Storing credentials in Kore")
	err = o.WaitForCreation(
		o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("ekscredential")).
			Name(o.Name).
			Payload(cred).
			Result(&eks.EKSCredentials{}),
		o.NoWait,
	)
	if err != nil {
		return fmt.Errorf("Error while creating credential: %s", err)
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
				Group:     eks.GroupVersion.Group,
				Version:   eks.GroupVersion.Version,
				Kind:      "EKSCredentials",
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
