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
	"os"
	"strings"

	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
)

// ConfigureOptions are the options of the command
type ConfigureOptions struct {
	cmdutil.Factory
	// Name is the name of the profile to configure
	Name string
	// Endpoint is the api endpoint to use
	Endpoint string
	// Force indicates we will overwrite any existing profiles
	Force bool
	// Account indicates the type of account
	Account string
	// LocalUser is an optional local username (used by basicauth)
	LocalUser string
	// LocalPass is an optional local password (user by basicauth)
	LocalPass string
}

// NewCmdProfilesConfigure creates and returns the profile configure command
func NewCmdProfilesConfigure(factory cmdutil.Factory) *cobra.Command {
	o := &ConfigureOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "configure a new profile for you",
		Example: "kore profile configure <name>",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := cmd.Flags()
	flags.BoolVarP(&o.Force, "force", "", false, "if true it overrides an existing profile with the same name")
	flags.StringVarP(&o.Endpoint, "api-url", "a", "", "url for the kore api endpoint `URL`")
	flags.StringVar(&o.Account, "account", "sso", "indicates the type of account for this profile `ACCOUNT`")
	flags.StringVarP(&o.LocalUser, "user", "u", "", "username when configuring basicauth profile `USERNAME`")
	flags.StringVar(&o.LocalPass, "password", "", "password for basicauth profiles ('-' for stdin) `PASSWORD`")

	return cmd
}

// Validate checks the options
func (o *ConfigureOptions) Validate() error {
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}
	if !utils.Contains(o.Account, kore.SupportedAccounts) {
		return errors.ErrUnknownAccountType
	}
	if o.Endpoint != "" && !utils.IsValidURL(o.Endpoint) {
		return fmt.Errorf("invalid api endpoint")
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

	if o.Endpoint == "" {
		prompts := cmdutil.Prompts{
			&cmdutil.Prompt{
				Id:    "Please enter the Kore API URL: (e.g https://api.domain.com)",
				Value: &o.Endpoint,
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
	}

	o.Endpoint = strings.TrimRight(o.Endpoint, "/")

	// @step: create an empty endpoint for then
	config.CreateProfile(o.Name, o.Endpoint)
	config.CurrentProfile = o.Name

	// @step: is the profile local?
	switch o.Account {
	case kore.AccountLocal, kore.AccountToken:
		if err := o.GetLocalAccountDetails(); err != nil {
			return err
		}
	}

	// @step: attempt to update the configuration
	if err := o.UpdateConfig(); err != nil {
		return fmt.Errorf("trying to update your local $ kore profile configure <name>: %s", err)
	}
	o.Println("Successfully configured the profile to: %s", o.Name)

	if o.Account == kore.AccountSSO {
		o.Println("Authenticate by running: $ kore login")
	}

	return nil
}

// GetLocalAccountDetails retrieves the local account settings
func (o *ConfigureOptions) GetLocalAccountDetails() error {
	switch o.Account {
	case kore.AccountLocal:
		auth := &config.BasicAuth{Username: o.LocalUser}

		if o.LocalUser == "" {
			p := cmdutil.Prompts{
				&cmdutil.Prompt{
					Id:     "Please enter your username",
					Value:  &auth.Username,
					ErrMsg: "invalid username",
				},
			}
			if err := p.Collect(); err != nil {
				return err
			}
		}

		if o.LocalPass == "" {
			p := cmdutil.Prompts{
				&cmdutil.Prompt{
					Id:     "Please enter your password",
					Value:  &auth.Password,
					ErrMsg: "invalid password",
				},
			}
			if err := p.Collect(); err != nil {
				return err
			}
		}
		if o.LocalPass != "" && o.LocalPass == "-" {
			pass, err := utils.ReadFileOrStdin(os.Stdin, o.LocalPass)
			if err != nil {
				return err
			}
			o.LocalPass = string(pass)
		}

		o.Config().AddAuthInfo(o.Name, &config.AuthInfo{BasicAuth: auth})

	case kore.AccountToken:
		var token string
		prompts := cmdutil.Prompts{
			&cmdutil.Prompt{Id: "Please enter your api token", Value: &token, ErrMsg: "invalid api token"},
		}
		if err := prompts.Collect(); err != nil {
			return err
		}
		o.Config().AddAuthInfo(o.Name, &config.AuthInfo{Token: &token})
	}

	return nil
}
