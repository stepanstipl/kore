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
	"github.com/appvia/kore/pkg/kore"

	"github.com/spf13/cobra"
)

// NewCmdDeleteIdentity creates and returns the create identity command
func NewCmdDeleteIdentity(factory cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "identity",
		Short:   "Deletes or one more identities in the kore",
		Example: "kore delete identity <type> [options]",
		Run:     cmdutil.RunHelp,
	}

	cmd.AddCommand(
		NewCmdDeleteIDPIdentity(factory),
		NewCmdDeleteBasicAuthIdentity(factory),
	)

	return cmd
}

// DeleteIDPOptions are the options for the command
type DeleteIDPOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the username
	Name string
}

// NewCmdDeleteIDPIdentity returns the delete sso command
func NewCmdDeleteIDPIdentity(factory cmdutil.Factory) *cobra.Command {
	o := &DeleteIDPOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     kore.AccountSSO,
		Short:   "Deletes a external identity from a kore managed user",
		Example: "kore delete identity sso [name]",
		Run:     cmdutil.DefaultRunFunc(o),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			suggestions, err := o.Resources().LookupResourceNames("user", "")
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return suggestions, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// Run implements the action
func (o *DeleteIDPOptions) Run() error {
	if o.Name == "" {
		who, err := o.Whoami()
		if err != nil {
			return err
		}
		o.Name = who.Username
	}

	err := o.ClientWithResource(o.Resources().MustLookup("identity")).
		Name(o.Name).
		SubResource(kore.AccountSSO).
		Delete().
		Error()
	if err != nil {
		return err
	}
	o.Println("Successfully removed the %s identity for %s", kore.AccountSSO, o.Name)

	return nil
}

// DeleteBasicOptions are the options for the command
type DeleteBasicOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Name is the username
	Name string
}

// NewCmdDeleteBasicAuthIdentity returns the delete basicauth command
func NewCmdDeleteBasicAuthIdentity(factory cmdutil.Factory) *cobra.Command {
	o := &DeleteBasicOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     kore.AccountLocal,
		Short:   "Deletes in locally managed identity in kore",
		Example: "kore delete identity basicauth [name]",
		Run:     cmdutil.DefaultRunFunc(o),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			suggestions, err := o.Resources().LookupResourceNames("user", "")
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return suggestions, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// Run implements the action
func (o *DeleteBasicOptions) Run() error {
	if o.Name == "" {
		who, err := o.Whoami()
		if err != nil {
			return err
		}
		o.Name = who.Username
	}

	err := o.ClientWithResource(o.Resources().MustLookup("identity")).
		Name(o.Name).
		SubResource(kore.AccountLocal).
		Delete().
		Error()
	if err != nil {
		return err
	}
	o.Println("Successfully removed the %s identity for %s", kore.AccountLocal, o.Name)

	return nil
}
