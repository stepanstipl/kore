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
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	longAssociateDesciption = `
Kore presently supports the use of multiple identities i.e. a user
can login via single sign, api token or basicauth (depending on kore's
configuration).

The associate command allows the user to associate an external idp
user configured to work along side Kore to a user already in Kore i.e.
assumed you've got a local user account and now added an SSO, you can
associate that account with local user.

Note, you must have your local credentials already setup in your
local kore configuration.
`

	longAssociateExamples = `
# Create an associate between local user and idp provider
$ kore create identity associate
`
)

// AssociateOptions are the options for the command
type AssociateOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
}

// NewCmdCreateAssociationIdentity creates and returns the command
func NewCmdCreateAssociationIdentity(factory cmdutil.Factory) *cobra.Command {
	o := &AssociateOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     "associate",
		Long:    longAssociateDesciption,
		Short:   "Create an associate between your local user and idp",
		Example: longAssociateExamples,
		Hidden:  true,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	return cmd
}

// Run implements the action
func (o *AssociateOptions) Run() error {
	return nil
}
