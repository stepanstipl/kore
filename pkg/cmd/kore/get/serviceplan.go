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

	"github.com/appvia/kore/pkg/client"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetServicePlanOptions the are the options for a get serviceplan ommand
type GetServicePlanOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is an optional name for the resource
	Name string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
	// Kind filters the service plans for a specific service kind
	Kind string
}

// NewCmdGetServicePlan creates and returns the get serviceplan command
func NewCmdGetServicePlan(factory cmdutil.Factory) *cobra.Command {
	o := &GetServicePlanOptions{Factory: factory}

	resource := o.Resources().MustLookup("serviceplan")

	command := &cobra.Command{
		Use:     "serviceplan",
		Aliases: []string{"serviceplans", resource.ShortName},
		Short:   "Returns all the service plans",
		Example: "kore get serviceplan [NAME] [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Kind, "kind", "k", "", "if set the command only returns the service plans with the given service kind")

	return command
}

// Validate is used to validate the options
func (o *GetServicePlanOptions) Validate() error {
	if o.Name != "" && o.Kind != "" {
		return errors.New("the --kind parameter should only be used when listing service plans")
	}

	return nil
}

// Run implements the action
func (o *GetServicePlanOptions) Run() error {
	resource := o.Resources().MustLookup("serviceplan")
	request := o.ClientWithResource(resource)

	if o.Name != "" {
		request.Name(o.Name)
	}

	if o.Kind != "" {
		request.Parameters(client.QueryParameter("kind", o.Kind))
	}

	if err := request.Get().Error(); err != nil {
		return err
	}

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
