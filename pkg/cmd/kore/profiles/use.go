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
	"errors"
	"fmt"

	cmderr "github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

type UseOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the name of the profile to use
	Name string
}

// NewCmdProfilesUse creates and returns the profile use command
func NewCmdProfilesUse(factory cmdutil.Factory) *cobra.Command {
	o := &UseOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "use",
		Short:   "switches to another profile",
		Example: "kore profile use <name>",

		PreRunE: cmdutil.RequireName,
		Run:     cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
			return o.Config().ListProfiles(), cobra.ShellCompDirectiveNoFileComp
		},
	}
	return command
}

// Run implements the actions
func (o *UseOptions) Run() error {
	config := o.Config()

	if !config.HasProfile(o.Name) {
		return errors.New("the profile does not exist")
	}
	if err := config.HasValidProfile(o.Name); err != nil {
		return cmderr.NewProfileInvalidError(err.Error(), o.Name)
	}
	config.CurrentProfile = o.Name

	if err := o.UpdateConfig(); err != nil {
		return fmt.Errorf("trying to update your local kore config: %s", err)
	}

	o.Println("Successfully switched the profile to: %s", o.Name)

	return nil
}
