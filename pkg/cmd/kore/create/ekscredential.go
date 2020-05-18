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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	confv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateEKSCredentialsOptions is used to provision a team
type CreateEKSCredentialsOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the name of the credential
	Name string
	// Description is a description of the credential
	Description string
	// DryRun indicates we only dryrun the resources
	DryRun bool
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

	flags := command.Flags()
	flags.StringVarP(&o.Description, "description", "d", "", "the description of the credential")
	flags.StringVarP(&o.AccountID, "account-id", "i", "", "the AWS numerical account ID")
	flags.StringVarP(&o.AccessKeyID, "key-id", "k", "", "the AWS access key ID")
	flags.StringVarP(&o.SecretAccessKey, "secret-key", "s", "", "the AWS secret access key - will prompt if not provided")
	flags.StringArrayVarP(&o.AllocateToTeams, "allocate", "a", []string{}, "list of teams to allocate to, e.g. team1,team2")
	flags.BoolVar(&o.AllocateToAll, "all-teams", false, "make these credentials available to all teams in kore (if not set, you must create an allocation for these credentials for them to be usable)")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutil.MustMarkFlagRequired(command, "account-id")
	cmdutil.MustMarkFlagRequired(command, "key-id")

	return command
}

// Run is responsible for creating the credentials
func (o CreateEKSCredentialsOptions) Run() error {
	var resources []runtime.Object

	secret, err := o.GenerateSecret()
	if err != nil {
		return err
	}
	resources = append(resources, secret)
	resources = append(resources, o.GenerateCredentials())
	resources = append(resources, o.GenerateAllocation())

	if o.DryRun {
		return cmdutil.RenderRuntimeObjectToYAML(resources, o.Writer())
	}

	for i := 0; i < len(resources); i++ {
		resource := resources[i]
		kind := strings.ToLower(kubernetes.GetRuntimeKind(resource))
		name := kubernetes.GetRuntimeName(resource)

		err := o.WaitForCreation(
			o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup(kind)).
				Name(name).
				Payload(resource),
			o.NoWait,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateSecret pulls in and generates the secret
func (o CreateEKSCredentialsOptions) GenerateSecret() (*configv1.Secret, error) {
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
			return nil, err
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
			Type:        "aws-credentials",
			Description: o.Description,
			Data: map[string]string{
				"access_key_id":     o.AccessKeyID,
				"access_secret_key": o.SecretAccessKey,
			},
		},
	}
	secret.Encode()

	return secret, nil
}

// GenerateCredentials is responsible for producing a gkecredentials
func (o *CreateEKSCredentialsOptions) GenerateCredentials() *eks.EKSCredentials {
	return &eks.EKSCredentials{
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
}

// GenerateAllocation is responsible for generating an allocation
func (o *CreateEKSCredentialsOptions) GenerateAllocation() *configv1.Allocation {
	teams := o.AllocateToTeams
	if o.AllocateToAll {
		teams = []string{"*"}
	}

	return &confv1.Allocation{
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
}
