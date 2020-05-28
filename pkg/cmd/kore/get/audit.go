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
	"time"

	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetAuditOptions the are the options for a get command
type GetAuditOptions struct {
	cmdutil.Factory
	// All indicates we should retrieve across all teams
	All bool
	// Headers indicates no headers on the table output
	Headers bool
	// Output is the output format
	Output string
	// Since is a time period to retrieve the events within
	Since time.Duration
	// Team string
	Team string
}

// NewCmdGetAudit creates and returns the get audit command
func NewCmdGetAudit(factory cmdutil.Factory) *cobra.Command {
	o := &GetAuditOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "audit",
		Short:   "Returns the audit log for teams or across kore",
		Example: "kore get audit [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.BoolVar(&o.All, "all", false, "retrieves logs across all teams (requires admin)")
	flags.DurationVarP(&o.Since, "since", "s", 30*time.Minute, "retrieves logs with the following period")

	return command
}

// Validate is responsible for checking the options
func (o *GetAuditOptions) Validate() error {
	if !o.All && o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// Run implements the action
func (o *GetAuditOptions) Run() error {
	resource := o.Resources().MustLookup("audit")
	request := o.ClientWithResource(resource).
		Parameters(
			client.QueryParameter("since", o.Since.String()),
		)

	if !o.All {
		request.Team(o.Team)
	}
	if err := request.Get().Error(); err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(request.Body())).
		ShowHeaders(o.Headers).
		Format(o.Output).
		Foreach("items").
		Printer(cmdutil.ConvertColumnsToRender(resource.Printer)...).
		Do()
}
