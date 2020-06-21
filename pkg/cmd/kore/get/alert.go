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
	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetAlertOptions the are the options for a get command
type GetAlertOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Team string
	Team string
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
	// AllTeams indicates we retrieve across all teams
	AllTeams bool
	// AllAlerts indicates we retrieve all alerts
	AllAlerts bool
}

// NewCmdGetAlert creates and returns the get admin command
func NewCmdGetAlert(factory cmdutil.Factory) *cobra.Command {
	o := &GetAlertOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Returns all the alerts held on resources in kore",
		Example: "kore get alerts [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.BoolVar(&o.AllTeams, "all-teams", false, "return the alerts across all teams (must be admin)")
	flags.BoolVarP(&o.AllAlerts, "all", "a", false, "returns all alerts event those ok")

	return command
}

// Validate is called to check the options
func (o *GetAlertOptions) Validate() error {
	if !o.AllTeams && o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// Run implements the action
func (o *GetAlertOptions) Run() error {
	var resp client.RestInterface

	resource := o.Resources().MustLookup("alert")
	params := []client.ParameterFunc{}
	if !o.AllAlerts {
		params = append(params,
			[]client.ParameterFunc{
				client.QueryParameter("status", "Active"),
				client.QueryParameter("status", "Silenced"),
			}...,
		)
	}

	if o.AllTeams {
		resp = o.ClientWithEndpoint("/monitoring/rules").Parameters(params...)
	} else {
		params = append(params, client.PathParameter("team", o.Team))
		resp = o.ClientWithEndpoint("/monitoring/teams/{team}/rules").Parameters(params...)
	}
	if err := resp.Get().Error(); err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(resp.Body())).
		ShowHeaders(o.Headers).
		Format(o.Output).
		Foreach("items").
		Printer(cmdutil.ConvertColumnsToRender(resource.Printer)...).
		Do()
}
