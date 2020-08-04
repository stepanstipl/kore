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
	"errors"
	"io/ioutil"
	"os"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	longBasicauthDesciption = `
Kore presently supports the use of multiple identities i.e. a user
can login via single sign, api token or basicauth (depending on kore's
configuration).

Note: for local user an administrator must first create and set the
password for the user.
`

	longBasicauthExamples = `
# Create a basicauth identity for the user james
$ kore create identity basicauth -u james -p - (password will be read from stdin)

# Update the current logged in user
$ kore create identity basicauth

# Users can update their identity via the same command
$ kore create identity basicauth -p - (defaults to current user)
`
)

// BasicAuthOptions are the options for the command
type BasicAuthOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Username is the user to create the identity for
	Username string
	// PassStdin indicates the password will come from stdin
	PassStdin bool
	// Password is the user password
	Password string
}

// NewCmdCreateBasicAuthIdentity creates and returns the command
func NewCmdCreateBasicAuthIdentity(factory cmdutil.Factory) *cobra.Command {
	o := &BasicAuthOptions{Factory: factory}

	cmd := &cobra.Command{
		Use:     "basicauth",
		Long:    longBasicauthDesciption,
		Short:   "Create a basicauth identity in kore",
		Example: longBasicauthExamples,
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := cmd.Flags()
	flags.StringVarP(&o.Username, "username", "u", "", "username you are updating the identity for (defaults to current) `USERNAME`")
	flags.BoolVarP(&o.PassStdin, "password", "p", false, "read the password from stdin `BOOL`")

	// @step: register the autocompletions
	cmdutil.MustRegisterFlagCompletionFunc(cmd, "username", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		suggestions, err := o.Resources().LookupResourceNames("user", "")
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

// Run implements the action
func (o *BasicAuthOptions) Run() error {
	// @step: if no username, lets default to current user
	who, err := o.Whoami()
	if err != nil {
		return err
	}
	if o.Username == "" {
		o.Username = who.Username
	}

	if o.PassStdin {
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		o.Password = string(content)
	} else {
		if err := (cmdutil.Prompts{
			&cmdutil.Prompt{
				Id:     "Please enter the password for " + o.Username,
				Value:  &o.Password,
				Mask:   true,
				ErrMsg: "invalid password",
			},
		}).Collect(); err != nil {
			return err
		}
		var confirm string

		if !o.PassStdin {
			if err := (cmdutil.Prompts{
				&cmdutil.Prompt{
					Id:     "Please confirm password for " + o.Username,
					Value:  &confirm,
					Mask:   true,
					ErrMsg: "invalid password",
				},
			}).Collect(); err != nil {
				return err
			}
			if confirm != o.Password {
				return errors.New("passwords do not match")
			}
		}
	}

	update := &orgv1.UpdateBasicAuthIdentity{
		Username: o.Username,
		Password: o.Password,
	}

	err = o.ClientWithResource(o.Resources().MustLookup("identity")).
		Name(o.Username).
		SubResource("basicauth").
		Payload(update).
		Update().
		Error()
	if err != nil {
		return err
	}
	o.Println("Successfully updated the identity in kore")

	// @step: if we have updated ourself we should update the config
	if o.Username == who.Username {
		auth := o.Config().GetAuthInfo(o.Client().CurrentProfile())
		if auth.BasicAuth != nil {
			auth.BasicAuth.Password = o.Password
			if err := o.UpdateConfig(); err != nil {
				return err
			}
		}
	}

	return nil
}
