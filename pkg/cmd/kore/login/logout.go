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

package login

import (
	"fmt"

	restconfig "github.com/appvia/kore/pkg/client/config"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

type LogoutOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is an optional profile name to logout from, else we use the current profile
	Name string
}

// NewCmdLogout is used to login to the api server
func NewCmdLogout(factory cmdutil.Factory) *cobra.Command {
	o := &LogoutOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "logout",
		Short:   "Deletes the login credentials from the current or selected profile",
		Example: "kore logout [name]",
		Run:     cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.BashCompDirective) {
			return o.Config().ListProfiles(), cobra.BashCompDirectiveNoFileComp
		},
	}

	return command
}

// Run implements the command
func (o *LogoutOptions) Run() error {
	name := o.Config().CurrentProfile
	if o.Name != "" {
		name = o.Name
	}

	profile, ok := o.Config().Profiles[name]
	if !ok {
		return fmt.Errorf("%q profile does not exist", name)
	}

	authInfo, ok := o.Config().AuthInfos[profile.AuthInfo]

	if !ok {
		return fmt.Errorf("%q user does not exist in configuration", profile.AuthInfo)
	}

	authInfo.OIDC = &restconfig.OIDC{
		ClientID:     authInfo.OIDC.ClientID,
		ClientSecret: authInfo.OIDC.ClientSecret,
		AuthorizeURL: authInfo.OIDC.AuthorizeURL,
	}

	if err := o.UpdateConfig(); err != nil {
		return err
	}
	o.Println("Logout successful")

	return nil
}
