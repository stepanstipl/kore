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
	"time"

	"github.com/appvia/kore/pkg/utils/render"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createServiceCredentialsLongDescription = `
Provisions service credentials for the given service and saves them as a Kubernetes secret in the target cluster
and namespace.

To list the available services
$ kore get services

To list the available clusters
$ kore get clusters

To list the available namespaces
$ kore get namespaceclaims

Examples:
$ kore -t <myteam> create servicecredentials db-creds --service my-database --cluster my-cluster --cluster-namespace dev

# Check the status of the service credentials
$ kore -t <myteam> get servicecredentials db-creds -o yaml
`
)

// CreateServiceCredentialsOptions is used to create service credentials
type CreateServiceCredentialsOptions struct {
	cmdutil.Factory
	// Name is the name of the service
	Name string
	// Credentials is the credentials allocation to build the cluster off
	Cluster string
	// Plan is the plan to build the service off
	Service string
	// Namespace is the target namespace in the cluster
	Namespace string
	// Team string
	Team string
	// Params is a collection of configuration parameters
	Params []string
	// NoWait indicates if we should wait for a service to provision
	NoWait bool
	// ShowTime indicate we should show the build time
	ShowTime bool
	// DryRun indicates that we should print the object instead of creating it
	DryRun bool
}

// NewCmdCreateServiceCredentials returns the create service credentials command
func NewCmdCreateServiceCredentials(factory cmdutil.Factory) *cobra.Command {
	o := &CreateServiceCredentialsOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "servicecredentials",
		Short:   "Creates service credentials within the team",
		Long:    createServiceCredentialsLongDescription,
		Example: "kore create servicecredentials -s <service> -c <cluster> -n <cluster namespace> [-t|--team]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Service, "service", "s", "", "service name `SERVICE`")
	flags.StringVarP(&o.Cluster, "cluster", "c", "", "cluster name `CLUSTER`")
	flags.StringVarP(&o.Namespace, "namespace", "n", "", "target namespace in the cluster `NAMESPACE`")
	flags.StringArrayVar(&o.Params, "param", []string{}, "a series of key value pairs used to override configuration parameters `KEY=VALUE`")
	flags.BoolVarP(&o.ShowTime, "show-time", "T", false, "shows the time it took to successfully provision a new service `BOOL`")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutils.MustMarkFlagRequired(command, "service")
	cmdutils.MustMarkFlagRequired(command, "cluster")
	cmdutils.MustMarkFlagRequired(command, "namespace")

	return command
}

// Validate is called to check the options
func (o *CreateServiceCredentialsOptions) Validate() error {
	if o.Team == "" {
		return errors.ErrTeamMissing
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	found, err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("service")).Name(o.Service).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(o.Service, "service")
	}

	found, err = o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("cluster")).Name(o.Cluster).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(o.Cluster, "cluster")
	}

	namespaceClaim := o.Cluster + "-" + o.Namespace

	found, err = o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("namespaceclaim")).Name(namespaceClaim).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(namespaceClaim, "namespaceclaim")
	}

	return nil
}

// Run implements the action
func (o *CreateServiceCredentialsOptions) Run() error {
	serviceCreds, err := o.CreateServiceCreds()
	if err != nil {
		return err
	}

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Resource(render.FromStruct(serviceCreds)).
			Format(render.FormatYAML).
			Do()
	}

	now := time.Now()
	// @step: provision and wait if required
	err = o.WaitForCreation(
		o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("servicecredentials")).
			Name(o.Name).
			Payload(serviceCreds),
		o.NoWait,
	)
	if err != nil {
		return err
	}
	if o.ShowTime {
		o.Println("Provisioning took: %s", time.Since(now))
	}

	return nil
}

// CreateServiceConfiguration is responsible for generating the service config
func (o *CreateServiceCredentialsOptions) CreateServiceCreds() (*servicesv1.ServiceCredentials, error) {
	service, err := o.GetService()
	if err != nil {
		return nil, err
	}

	cluster, err := o.GetCluster()
	if err != nil {
		return nil, err
	}

	configJSON, err := cmdutil.PatchJSON("{}", o.Params)
	if err != nil {
		return nil, err
	}

	serviceCreds := &servicesv1.ServiceCredentials{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceCredentials",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Team,
		},
		Spec: servicesv1.ServiceCredentialsSpec{
			Kind:             service.Spec.Kind,
			Service:          service.Ownership(),
			Cluster:          cluster.Ownership(),
			ClusterNamespace: o.Namespace,
			Configuration:    &apiextv1.JSON{Raw: []byte(configJSON)},
		},
	}

	return serviceCreds, nil
}

// GetService retrieves the requested service
func (o *CreateServiceCredentialsOptions) GetService() (*servicesv1.Service, error) {
	service := &servicesv1.Service{}
	return service, o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("service")).
		Name(o.Service).
		Result(service).
		Get().
		Error()
}

// GetService retrieves the requested service
func (o *CreateServiceCredentialsOptions) GetCluster() (*clustersv1.Cluster, error) {
	cluster := &clustersv1.Cluster{}
	return cluster, o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("cluster")).
		Name(o.Cluster).
		Result(cluster).
		Get().
		Error()
}
