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

package profiles

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var longProfileDescription = `
Profiles provide a means to store, configure and switch between multiple
Appvia Kore instances from a single configuration. Alternatively, you might
use profiles to use different identities (i.e. admin / user) to a single
instance. These are automatically created for you via the $ kore login
command or you can manually configure them via the $ kore profile configure.

Examples:

$ kore profile                     # will show this help menu
$ kore profile show                # will show the profile in use
$ kore profile list                # will show all the profiles available to you
$ kore profile use <name>          # switches to another profile
$ kore profile configure <name>    # allows you to configure a profile
$ kore profile rm <name>           # removes a profile
$ kore profile set <path> <value>  # set configuration values
`

// NewCmdProfiles creates and returns the profiles command
func NewCmdProfiles(factory cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "profiles",
		Aliases:               []string{"profile"},
		DisableFlagsInUseLine: true,
		Long:                  longProfileDescription,
		Short:                 "Manage profiles, allowing you switch, list and show profiles",
		Run:                   cmdutil.RunHelp,
	}

	cmd.AddCommand(
		NewCmdProfilesList(factory),
		NewCmdProfilesUse(factory),
		NewCmdProfilesShow(factory),
		NewCmdProfilesDelete(factory),
		NewCmdProfilesConfigure(factory),
		NewCmdProfilesSet(factory),
	)

	return cmd
}
