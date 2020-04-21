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
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetAdminOptions the are the options for a get command
type GetAdminOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
}

// NewCmdGetAdmin creates and returns the get admin command
func NewCmdGetAdmin(factory cmdutil.Factory) *cobra.Command {
	o := &GetAdminOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "admin",
		Short:   "Returns all the user in kore whom are adminstrators",
		Example: "kore get admin [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the action
func (o *GetAdminOptions) Run() error {
	client := o.ClientWithTeamResource(kore.HubAdminTeam, o.Resources().MustLookup("member"))
	if err := client.Get().Error(); err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(client.Body())).
		ShowHeaders(o.Headers).
		Format(o.Output).
		Foreach("items").
		Printer(
			render.Column("Username", "."),
		).Do()
}
