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

package get

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/appvia/kore/pkg/client"
	cmderrors "github.com/appvia/kore/pkg/cmd/errors"
	"github.com/appvia/kore/pkg/cmd/kore/identity"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

var (
	getLongDescription = `
Allows to you retrieve the resources from the kore api. The command format
is <resource> [name]. When the optional name is not provided we will return
a full listing of all the <resource>s from the API. Examples of resource types
are users, teams, gkes, clusters amongst a few.

You can list all the available resource via $ kore api-resources

Though for a better experience all the resource are autocompletes for you.
Take a look at $ kore completion for details
`
	getExamples = `
# List users:
$ kore get users

#Get information about a specific user:
$ kore get user admin [-o yaml]
`
)

// GetOptions the are the options for a get command
type GetOptions struct {
	cmdutil.Factory
	// Name is an optional name for the resource
	Name string
	// Resource is the resource to retrieve
	Resource string
	// Team string
	Team string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
	// Raw URI to request from the server
	Raw string
}

// NewCmdGet creates and returns the get command
func NewCmdGet(factory cmdutil.Factory) *cobra.Command {
	o := &GetOptions{Factory: factory}

	// @step: retrieve a list of known resources
	possible, _ := factory.Resources().Names()

	command := &cobra.Command{
		Use:     "get",
		Long:    getLongDescription,
		Example: getExamples,
		Run:     cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				return possible, cobra.ShellCompDirectiveNoFileComp
			case 1:
				suggestions, err := o.Resources().LookupResourceNames(cmd.Flags().Arg(0), cmdutil.GetTeam(cmd))
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				// choice we don't want to show everything here
				if len(suggestions) > 15 {
					return suggestions[0:15], cobra.ShellCompDirectiveNoFileComp
				}

				return suggestions, cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	flags := command.Flags()
	flags.StringVar(&o.Raw, "raw", "", "Raw URI to request from the server")

	command.AddCommand(
		NewCmdGetAdmin(factory),
		NewCmdGetAlert(factory),
		NewCmdGetAudit(factory),
		identity.NewCmdGetIdentity(factory),
	)

	if factory.Config().FeatureGates[kore.FeatureGateServices] {
		command.AddCommand(
			NewCmdGetServiceKind(factory),
			NewCmdGetServicePlan(factory),
			NewCmdGetServiceCredential(factory),
		)
	}

	return command
}

// Validate is used to validate the options
func (o *GetOptions) Validate() error {
	if o.Resource == "" && o.Raw == "" {
		return cmderrors.ErrMissingResource
	}

	if o.Raw != "" && (o.Resource != "" || o.Name != "") {
		return errors.New("when --raw is set, no arguments should be used")
	}

	return nil
}

// Run implements the action
func (o *GetOptions) Run() error {
	if o.Raw != "" {
		return o.runRaw()
	}

	// @step: lookup the resource from the cache
	resource, err := o.Resources().Lookup(o.Resource)
	if err != nil {
		return err
	}

	// @step: if the resource if team space, lets ensure we have the team selector
	if resource.IsTeamScoped() && o.Team == "" {
		return cmderrors.ErrTeamMissing
	}

	// @step: we need to construct the request
	request := o.ClientWithResource(resource)

	if resource.IsScoped(cmdutil.TeamScope) {
		request.Team(o.Team)
	}
	if resource.IsScoped(cmdutil.DualScope) && o.Team != "" {
		request.Team(o.Team)
	}
	if o.Name != "" {
		request.Name(o.Name)
	}

	// @step: we perform the get request against the api
	if err := request.Get().Error(); err != nil {
		return err
	}

	// @step: construct the columns from the resource - this could probably be
	// cleaned up some how
	display := render.Render().
		Writer(o.Writer()).
		ShowHeaders(o.Headers).
		Format(o.Output).
		Resource(
			render.FromReader(request.Body()),
		).
		Printer(cmdutil.ConvertColumnsToRender(resource.Printer)...)

	if o.Name == "" {
		display.Foreach("items")
	}

	return display.Do()
}

func (o *GetOptions) runRaw() error {
	if o.Output != render.FormatJSON && o.Output != render.FormatYAML {
		o.Output = render.FormatJSON
	}

	u, err := url.Parse(o.Raw)
	if err != nil {
		return fmt.Errorf("raw url is invalid")
	}

	var params []client.ParameterFunc
	for key, values := range u.Query() {
		for _, value := range values {
			params = append(params, client.QueryParameter(key, value))
		}
	}

	resp := o.ClientWithEndpoint(u.Path).Parameters(params...).Get()
	if err := resp.Error(); err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(resp.Body())).
		Format(o.Output).
		Do()
}
