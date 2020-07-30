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

package identity

import (
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/client"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// GetIdentityOptions the are the options for a get command
type GetIdentityOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// All indicates we should retrieve all identities not just our own
	All bool
	// Output is the output format
	Output string
	// Headers indicates no headers on the table output
	Headers bool
}

// NewCmdGetIdentity creates and returns the get admin command
func NewCmdGetIdentity(factory cmdutil.Factory) *cobra.Command {
	o := &GetIdentityOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     "identity",
		Aliases: []string{"identities"},
		Short:   "Used to query user identities in kore",
		Example: "kore get identities [options]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.All, "all", false, "retrieve all identities, not just your own `BOOL`")

	return cmd
}

// Run implements the action
func (o *GetIdentityOptions) Run() error {
	user, err := o.Whoami()
	if err != nil {
		return err
	}

	list := &orgv1.IdentityList{}

	if o.All {
		err = o.ClientWithResource(o.Resources().MustLookup("identity")).
			Parameters(client.QueryParameter("all", "true")).
			Result(list).
			Get().
			Error()
	} else {
		err = o.ClientWithResource(o.Resources().MustLookup("identity")).
			Parameters(client.QueryParameter("all", "false")).
			Name(user.Username).
			Result(list).
			Get().
			Error()
	}
	if err != nil {
		return err
	}

	return render.Render().
		Writer(o.Writer()).
		ShowHeaders(o.Headers).
		Format(o.Output).
		Foreach("items").
		Resource(
			render.FromStruct(list),
		).
		Printer(
			render.Column("Username", "spec.user.metadata.name"),
			render.Column("Account Type", "spec.accountType"),
		).
		Do()
}
