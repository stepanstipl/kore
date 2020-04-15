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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createNamespaceLongDescription = `
Provides the ability to create a namespace on a provisioned cluster. In order to
retrieve the clusters you have available you can run:

$ kore get clusters -t <team>

Examples:
# Create a namespace on cluster 'dev'
$ kore create namespace -c cluster -t <team>

# Deleting a namespace on the cluster
$ kore delete namespaceclaim

You can list the namespace you have already provisioned via

$ kore get namespaceclaims -t <team>
`
)

// NamespaceOptions is used to provision a team
type NamespaceOptions struct {
	cmdutil.Factory
	// Cluster is the cluster you are creating the namespace in
	Cluster string
	// Force is used to force an operation
	Force bool
	// Name is the name of the namespace
	Name string
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
	// Team is the team name
	Team string
}

// NewCmdCreateNamespace returns the create namespace command
func NewCmdCreateNamespace(factory cmdutil.Factory) *cobra.Command {
	o := &NamespaceOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "namespace",
		Short:   "Creates a namespace within a managed cluster",
		Long:    createNamespaceLongDescription,
		Example: "kore create namespace -u <cluster> [-t|--team]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	command.Flags().StringVarP(&o.Cluster, "cluster", "c", "", "the name of the cluster you are creating the namespace on `NAME`")
	cmdutils.MustMarkFlagRequired(command, "cluster")

	// @step: add auto complete on the cluster name
	cmdutils.MustRegisterFlagCompletionFunc(command, "cluster", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.BashCompDirective) {
		suggestions, err := o.Resources().LookupResourceNames("cluster", cmdutil.GetTeam(cmd))
		if err != nil {
			return nil, cobra.BashCompDirectiveError
		}

		return suggestions, cobra.BashCompDirectiveNoFileComp
	})

	return command
}

// Validate is responsible for checking the options
func (o *NamespaceOptions) Validate() error {
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	found, err := o.Client().
		Team(o.Team).
		Resource("cluster").
		Name(o.Cluster).
		Exists()
	if err != nil {
		return err
	}

	if !found && !o.Force {
		return errors.NewResourceNotFound(o.Cluster)
	}

	return nil
}

// Run implements the action
func (o *NamespaceOptions) Run() error {
	name := fmt.Sprintf("%s-%s", o.Cluster, o.Name)

	namespace := &clustersv1.NamespaceClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "NamespaceClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: o.Team,
		},
		Spec: clustersv1.NamespaceClaimSpec{
			Name: o.Name,
			Cluster: corev1.Ownership{
				Group:     clustersv1.GroupVersion.Group,
				Version:   clustersv1.GroupVersion.Version,
				Kind:      "Kubernetes",
				Namespace: o.Team,
				Name:      o.Cluster,
			},
		},
	}

	return o.WaitForCreation(
		o.Client().
			Team(o.Team).
			Resource("namespaceclaim").
			Payload(namespace).
			Name(name),
		o.NoWait,
	)
}
