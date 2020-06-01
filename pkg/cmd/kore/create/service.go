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
	"regexp"
	"time"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
	apiexts "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createServiceLongDescription = `
Provides the ability to provision a service in the team. The service
itself is created from a predefined service plan (a template). You can view the service plans
available to you via $ kore get serviceplans.
`
	createServiceExamples = `
Note: you can retrieve a list of all the service plans available to you via:

$ kore get serviceplans
$ kore get serviceplan <name> -o yaml

# Create a service foo from plan some-plan
$ kore -t <myteam> create service foo --plan some-plan

# You can override the plan parameters using the --param
$ kore -t <myteam> create service foo --param configkey=value

# You can using JSON values when setting a parameter
$ kore -t <myteam> create service foo --param 'configlist=[1, 2, 3]'

# Check the status of the service
$ kore -t <myteam> get service foo -o yaml
`
)

// CreateServiceOptions is used to create a service
type CreateServiceOptions struct {
	cmdutil.Factory
	// Name is the name of the service
	Name string
	// Description is a description of the service
	Description string
	// Plan is the plan to build the service off
	Plan string
	// Team string
	Team string
	// PlanParams is a collection of service plan configuration overrides
	PlanParams []string
	// NoWait indicates if we should wait for a service to provision
	NoWait bool
	// ShowTime indicate we should show the build time
	ShowTime bool
	// DryRun indicates that we should print the object instead of creating it
	DryRun bool
}

// NewCmdCreateService returns the create service command
func NewCmdCreateService(factory cmdutil.Factory) *cobra.Command {
	o := &CreateServiceOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "service",
		Short:   "Create a service within the team",
		Long:    createServiceLongDescription,
		Example: createServiceExamples,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Plan, "plan", "p", "", "plan which this service will be templated from `NAME`")
	flags.StringVarP(&o.Description, "description", "d", "", "a short description for the service `DESCRIPTION`")
	flags.StringArrayVar(&o.PlanParams, "param", []string{}, "a series of key value pairs used to override plan parameters  `KEY=VALUE`")
	flags.BoolVarP(&o.ShowTime, "show-time", "T", false, "shows the time it took to successfully provision a new service `BOOL`")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutil.MustMarkFlagRequired(command, "plan")

	cmdutil.MustRegisterFlagCompletionFunc(command, "plan", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		suggestions, err := o.Resources().LookupResourceNames("serviceplan", "")
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	return command
}

// Validate is called to check the options
func (o *CreateServiceOptions) Validate() error {
	if o.Team == "" {
		return errors.ErrTeamMissing
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	found, err := o.ClientWithResource(o.Resources().MustLookup("serviceplan")).Name(o.Plan).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(o.Plan, "serviceplan")
	}

	match := regexp.MustCompile("^.*=.*$")

	for _, x := range o.PlanParams {
		if !match.MatchString(x) {
			return errors.NewInvalidParamError("param", x)
		}
	}

	return nil
}

// Run implements the action
func (o *CreateServiceOptions) Run() error {
	service, err := o.CreateService()
	if err != nil {
		return err
	}

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Resource(render.FromStruct(service)).
			Format(render.FormatYAML).
			Do()
	}

	now := time.Now()
	// @step: provision and wait if required
	err = o.WaitForCreation(
		o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("service")).
			Name(o.Name).
			Payload(service),
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

// CreateService is responsible for generating the service config
func (o *CreateServiceOptions) CreateService() (*servicesv1.Service, error) {
	plan, err := o.GetPlan()
	if err != nil {
		return nil, err
	}

	var configJSON string
	if plan.Spec.Configuration != nil {
		configJSON = string(plan.Spec.Configuration.Raw)
	}

	if configJSON, err = cmdutil.PatchJSON(configJSON, o.PlanParams); err != nil {
		return nil, err
	}

	service := &servicesv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: servicesv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Team,
		},
		Spec: servicesv1.ServiceSpec{
			Kind:          plan.Spec.Kind,
			Plan:          plan.Name,
			Configuration: &apiexts.JSON{Raw: []byte(configJSON)},
		},
	}

	return service, nil
}

// GetPlan retrieves the requested service plan
func (o *CreateServiceOptions) GetPlan() (*servicesv1.ServicePlan, error) {
	plan := &servicesv1.ServicePlan{}

	return plan, o.ClientWithResource(o.Resources().MustLookup("serviceplan")).
		Name(o.Plan).
		Result(plan).
		Get().
		Error()
}
