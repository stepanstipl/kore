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

// NewCmdCreateIdentity creates and returns the create identity command
func NewCmdCreateIdentity(factory cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity",
		Short:   "Creates or one more identities in the kore",
		Example: "kore create identity <type> [options]",
		Run:     cmdutil.RunHelp,
	}

	cmd.AddCommand(
		NewCmdCreateAssociationIdentity(factory),
		NewCmdCreateBasicAuthIdentity(factory),
	)

	return cmd
}
