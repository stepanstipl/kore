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
	kerrors "github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetServiceCredentialOptions the are the options for a get servicecredential command
type GetServiceCredentialOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Team string
	Team string
	// Name is an optional name for the resource
	Name string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
	// Service filters the service credentials for a specific service
	Service string
	// Cluster filters the service credentials for a specific cluster
	Cluster string
}

// NewCmdGetServiceCredential creates and returns the get servicecredential command
func NewCmdGetServiceCredential(factory cmdutil.Factory) *cobra.Command {
	o := &GetServiceCredentialOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "servicecredential",
		Aliases: []string{"servicecredentials"},
		Short:   "Returns all the service plans",
		Example: "kore get servicecredential [NAME] [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Service, "service", "s", "", "if set the command only returns the service credentials for the given service")
	flags.StringVarP(&o.Cluster, "cluster", "c", "", "if set the command only returns the service credentials for the given cluster")

	return command
}

// Validate is used to validate the options
func (o *GetServiceCredentialOptions) Validate() error {
	if o.Team == "" {
		return kerrors.ErrTeamMissing
	}
	if o.Name != "" && o.Service != "" {
		return errors.New("the --service parameter should only be used when listing service credentials")
	}

	if o.Name != "" && o.Cluster != "" {
		return errors.New("the --cluster parameter should only be used when listing service credentials")
	}

	return nil
}

// Run implements the action
func (o *GetServiceCredentialOptions) Run() error {
	resource := o.Resources().MustLookup("servicecredential")
	request := o.ClientWithTeamResource(o.Team, resource)

	if o.Name != "" {
		request.Name(o.Name)
	}
	if o.Service != "" {
		request.Parameters(client.QueryParameter("service", o.Service))
	}
	if o.Cluster != "" {
		request.Parameters(client.QueryParameter("cluster", o.Cluster))
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
