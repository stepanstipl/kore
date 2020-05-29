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

	"github.com/appvia/kore/pkg/client"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetServiceKindOptions the are the options for a get servicekind ommand
type GetServiceKindOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	cmd *cobra.Command
	// Name is an optional name for the resource
	Name string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
	// Platform filters the service kinds for a specific service platform
	Platform string
	// Enabled filters for the enabled/disabled status
	Enabled bool
}

// NewCmdGetServiceKind creates and returns the get servicekind command
func NewCmdGetServiceKind(factory cmdutil.Factory) *cobra.Command {
	o := &GetServiceKindOptions{Factory: factory}

	resource := o.Resources().MustLookup("servicekind")

	command := &cobra.Command{
		Use:     "servicekind",
		Aliases: []string{"servicekinds", resource.ShortName},
		Short:   "Returns all the service plans",
		Example: "kore get servicekind [NAME] [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}
	o.cmd = command

	flags := command.Flags()
	flags.StringVarP(&o.Platform, "platform", "p", "", "if set the command only returns the service kinds for the given service platform")
	flags.BoolVarP(&o.Enabled, "enabled", "e", true, "if set the command only returns the enabled/disabled service kinds")

	return command
}

// Validate is used to validate the options
func (o *GetServiceKindOptions) Validate() error {
	if o.Name != "" && o.Platform != "" {
		return errors.New("the --platform parameter should only be used when listing service kinds")
	}

	if o.Name != "" && o.cmd.Flag("enabled").Changed {
		return errors.New("the --enabled parameter should only be used when listing service kinds")
	}

	return nil
}

// Run implements the action
func (o *GetServiceKindOptions) Run() error {
	resource := o.Resources().MustLookup("servicekind")
	request := o.ClientWithResource(resource)

	if o.Name != "" {
		request.Name(o.Name)
	}

	if o.Platform != "" {
		request.Parameters(client.QueryParameter("platform", o.Platform))
	}

	if o.cmd.Flag("enabled").Changed {
		request.Parameters(client.QueryParameter("enabled", fmt.Sprintf("%v", o.Enabled)))
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
