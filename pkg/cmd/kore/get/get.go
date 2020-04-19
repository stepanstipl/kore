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
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

var getLongDescription = `
Allows to you retrieve the resources from the kore api. The command format
is <resource> [name]. When the optional name is not provided we will return
a full listing of all the <resource>s from the API. Examples of resource types
are users, teams, gkes, clusters amongst a few.

You can list all the available resource via $ kore api-resources

Though for a better experience all the resource are autocompletes for you.
Take a look at $ kore completion for details

Examples:
List users:
$ korectl get users

Get information about a specific user:
$ korectl get user admin [-o yaml]
`

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
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
	// Headers indicates no headers on the table output
	Headers bool
}

// NewCmdGet creates and returns the get command
func NewCmdGet(factory cmdutil.Factory) *cobra.Command {
	o := &GetOptions{Factory: factory}

	// @step: retrieve a list of known resources
	possible, _ := factory.Resources().Names()

	command := &cobra.Command{
		Use:     "get",
		Short:   "Retrieve resources from the kore api",
		Long:    getLongDescription,
		Example: "kore get <resource> [name] [options]",
		Run:     cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				return append(possible, "admin"), cobra.ShellCompDirectiveNoFileComp
			case 1:
				resource := o.Resources().ResolveShorthand(cmd.Flags().Arg(0))
				suggestions, err := o.Resources().LookupResourceNames(resource, cmdutil.GetTeam(cmd))
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				// choice we don't wan't to show everything here
				if len(suggestions) > 15 {
					return suggestions[0:15], cobra.ShellCompDirectiveNoFileComp
				}

				return suggestions, cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	command.AddCommand(
		NewCmdGetAdmin(factory),
		NewCmdGetAudit(factory),
	)

	return command
}

// Validate is used to validate the options
func (o *GetOptions) Validate() error {
	if o.Resource == "" {
		return errors.ErrMissingResource
	}

	resource, err := o.Resources().Lookup(utils.Pluralize(o.Resource))
	if err != nil {
		return err
	}
	if resource == nil {
		return errors.ErrUnknownResource
	}

	return nil
}

// Run implements the action
func (o *GetOptions) Run() error {
	plural := utils.Pluralize(o.Resource)

	// @step: lookup the resource from the cache
	resource, err := o.Resources().Lookup(plural)
	if err != nil {
		return err
	}
	if resource == nil {
		return errors.ErrUnknownResource
	}

	// @step: if the resource if team space, lets ensure we have the team selector
	if resource.IsTeamScoped() && o.Team == "" {
		return errors.ErrTeamMissing
	}

	// @step: we need to construct the request
	request := o.Client().Resource(plural)

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
		Resource(render.FromReader(request.Body())).
		ShowHeaders(o.Headers).
		Format(o.Output)

	columns := make([]render.PrinterColumnFunc, len(resource.Printer))
	for i, c := range resource.Printer {
		switch c.Format {
		case "age":
			columns[i] = render.Column(c.Name, c.Path, render.Age())
		default:
			columns[i] = render.Column(c.Name, c.Path)
		}
	}
	display.Printer(columns...)

	if o.Name == "" {
		display.Foreach("items")
	}

	return display.Do()
}
