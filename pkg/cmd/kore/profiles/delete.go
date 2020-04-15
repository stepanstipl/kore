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
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cmdutil.Factory
	// Name is the profile to delete
	Name string
}

// NewCmdProfilesDelete creates and returns the profile delete command
func NewCmdProfilesDelete(factory cmdutil.Factory) *cobra.Command {
	o := &DeleteOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "removes a profile from configuration",
		Example: "korectl profile delete <name>",
		Run:     cmdutil.DefaultRunFunc(o),

		ValidArgsFunction: func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.BashCompDirective) {
			return o.Config().ListProfiles(), cobra.BashCompDirectiveNoFileComp
		},
	}

	return command
}

// Validate checks the options
func (o DeleteOptions) Validate() error {
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	if !o.Config().HasProfile(o.Name) {
		return errors.ErrMissingProfile
	}

	return nil
}

// Run implements the action
func (o *DeleteOptions) Run() error {
	config := o.Config()

	if config.CurrentProfile == o.Name {
		config.CurrentProfile = ""
	}
	config.RemoveProfile(o.Name)

	if err := o.UpdateConfig(); err != nil {
		return err
	}
	o.Println("Successfully removed the profile: %s", o.Name)

	return nil

}
