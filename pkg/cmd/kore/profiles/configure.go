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
	"fmt"
	"strings"

	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
)

type ConfigureOptions struct {
	cmdutil.Factory
	// Name is the name of the profile to configure
	Name string
	// Force indicates we will overwrite any existing profiles
	Force bool
}

// NewCmdProfilesConfigure creates and returns the profile configure command
func NewCmdProfilesConfigure(factory cmdutil.Factory) *cobra.Command {
	o := &ConfigureOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "configure a new profile for you",
		Example: "kore profile configure <name>",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	command.Flags().BoolVarP(&o.Force, "force", "", false, "if true it overrides an existing profile with the same name")

	return command
}

// Validate checks the options
func (o *ConfigureOptions) Validate() error {
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	return nil
}

// Run implements the action
func (o *ConfigureOptions) Run() error {
	config := o.Config()

	// @check the profile does not conflict
	if !o.Force && config.HasProfile(o.Name) {
		return errors.NewConflictError("profile name is already taken, please choose another name")
	}

	var endpoint string
	prompts := cmdutil.Prompts{
		&cmdutil.Prompt{
			Id:    "Please enter the Kore API URL: (e.g https://api.domain.com)",
			Value: &endpoint,
			Validate: func(in string) error {
				if !utils.IsValidURL(in) {
					return fmt.Errorf("invalid endpoint: %s", in)
				}
				return nil
			},
		},
	}
	if err := prompts.Collect(); err != nil {
		return err
	}

	endpoint = strings.TrimRight(endpoint, "/")

	// @step: create an empty endpoint for then
	config.CreateProfile(o.Name, endpoint)
	config.CurrentProfile = o.Name

	// @step: attempt to update the configuration
	if err := o.UpdateConfig(); err != nil {
		return fmt.Errorf("trying to update your local $ kore profile configure <name>: %s", err)
	}

	o.Println("Successfully configured the profile to: %s", o.Name)
	o.Println("In order to authenticate please run: $ kore login")

	return nil

}
