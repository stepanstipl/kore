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
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/client"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	apiexts "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createServiceLongDescription = `
Provides the ability to provision a service in the team. The service
itself is created from a predefined service plan (a template). You can view the service plans
available to you via $ kore get serviceplans.

Note: you can retrieve a list of all the service plans available to you via:
$ kore get serviceplans
$ kore get serviceplan <name> -o yaml

Examples:
$ kore -t <myteam> create service foo --plan some-plan

# Check the status of the service
$ kore -t <myteam> get service foo -o yaml
`
)

// CreateServiceOptions is used to create a service
type CreateServiceOptions struct {
	cmdutil.Factory
	// Name is the name of the service
	Name string
	// Credentials is the credentials allocation to build the cluster off
	Credentials string
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
}

// NewCmdCreateService returns the create service command
func NewCmdCreateService(factory cmdutil.Factory) *cobra.Command {
	o := &CreateServiceOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "service",
		Short:   "Create a service within the team",
		Long:    createServiceLongDescription,
		Example: "kore create service -p <plan> [-t|--team]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Credentials, "credentials", "c", "", "name of the credentials allocation to use for this service `NAME`")
	flags.StringVarP(&o.Plan, "plan", "p", "", "plan which this service will be templated from `NAME`")
	flags.StringVarP(&o.Description, "description", "d", "", "a short description for the service `DESCRIPTION`")
	flags.StringSliceVar(&o.PlanParams, "param", []string{}, "preprovision a collection namespaces on this service as well `NAMES`")
	flags.BoolVarP(&o.ShowTime, "show-time", "T", false, "shows the time it took to successfully provision a new service `BOOL`")

	cmdutils.MustMarkFlagRequired(command, "plan")

	cmdutils.MustRegisterFlagCompletionFunc(command, "plan", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
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
	config, err := o.CreateServiceConfiguration()
	if err != nil {
		return err
	}

	now := time.Now()
	// @step: provision and wait if required
	err = o.WaitForCreation(
		o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("service")).
			Name(o.Name).
			Payload(config),
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
func (o *CreateServiceOptions) CreateServiceConfiguration() (*servicesv1.Service, error) {
	plan, err := o.GetPlan()
	if err != nil {
		return nil, err
	}

	var configuration map[string]interface{}

	if err := json.Unmarshal(plan.Spec.Configuration.Raw, &configuration); err != nil {
		return nil, fmt.Errorf("failed to parse plan configuration: %s", err)
	}

	params, err := o.ParsePlanParams()
	if err != nil {
		return nil, err
	}

	for k, v := range params {
		configuration[k] = v
	}

	cc := &bytes.Buffer{}
	if err := json.NewEncoder(cc).Encode(configuration); err != nil {
		return nil, fmt.Errorf("failed to process service configuration: %s", err)
	}

	credentials, err := o.GetCredentialsAllocation()
	if err != nil {
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
			Configuration: apiexts.JSON{Raw: cc.Bytes()},
			Credentials:   credentials.Spec.Resource,
		},
	}

	return service, nil
}

// ParsePlanParams is responsible for parsing the plan overrides
func (o *CreateServiceOptions) ParsePlanParams() (map[string]interface{}, error) {
	params := map[string]interface{}{}

	for _, param := range o.PlanParams {
		parts := regexp.MustCompile(`\s*=\s*`).Split(strings.TrimSpace(param), 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid plan-param value %q, you must use param=<JSON value> format", param)
		}
		name := parts[0]
		jsonValue := parts[1]

		var parsed interface{}
		if err := json.Unmarshal([]byte(jsonValue), &parsed); err != nil {
			return nil, err
		}
		params[name] = parsed
	}

	return params, nil
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

// GetAllocation retrieves the requested allocation
func (o *CreateServiceOptions) GetCredentialsAllocation() (*configv1.Allocation, error) {
	allocation := &configv1.Allocation{}

	if o.Credentials == "" {
		return allocation, nil
	}

	err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("allocation")).
		Name(o.Credentials).
		Result(allocation).
		Get().
		Error()
	if err != nil {
		if client.IsNotFound(err) {
			return nil, fmt.Errorf("credentials allocation %q does not exist", o.Credentials)
		}
		return nil, err
	}

	return allocation, nil
}
