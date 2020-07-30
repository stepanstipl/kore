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

package kore

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
)

// WhoAmIOptions are the options for the whoami command
type WhoAmIOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
}

// NewCmdWhoami returns the whoami command
func NewCmdWhoami(factory cmdutil.Factory) *cobra.Command {
	o := &WhoAmIOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "whoami",
		Aliases: []string{"who"},
		Short:   "Used to retrieve details on your identity within the kore",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the action
func (o *WhoAmIOptions) Run() error {
	resp := o.ClientWithEndpoint("/whoami").Get()
	if resp.Error() != nil {
		return resp.Error()
	}

	return render.Render().
		Writer(o.Writer()).
		Resource(render.FromReader(resp.Body())).
		Unknown("None").
		Printer(
			render.Column("Username", "username"),
			render.Column("Email", "email"),
			render.Column("Teams", "teams|@sjoin"),
			render.Column("Authentication", "authMethod"),
		).Do()
}
