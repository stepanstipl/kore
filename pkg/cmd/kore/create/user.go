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

package create

import (
	"fmt"
	"strings"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	userLongDesciption = `
Provides the ability to provision a user in kore.

# Create the test user
$ kore create username test -e test@appiva.io
`
)

// CreateUserOptions is used to provision a team
type CreateUserOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
	// Name is the username to add
	Name string
	// Email is the user email address
	Email string
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
}

// NewCmdCreateUser returns the create user command
func NewCmdCreateUser(factory cmdutils.Factory) *cobra.Command {
	o := &CreateUserOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "user",
		Short:   "Adds to the user to kore",
		Long:    userLongDesciption,
		Example: "kore create user <username> -e <email> [-t team]",
		PreRunE: cmdutils.RequireName,
		Run:     cmdutils.DefaultRunFunc(o),
	}

	command.Flags().StringVarP(&o.Email, "email", "e", "", "an email address for the user `EMAIL`")
	cmdutils.MustMarkFlagRequired(command, "email")

	// @step: register the autocompletions
	cmdutils.MustRegisterFlagCompletionFunc(command, "email", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		name := cmd.Flags().Arg(0)
		if name == "" {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		list := &orgv1.UserList{}
		if err := o.Client().Resource("user").Result(list).Get().Error(); err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var domains []string
		for _, x := range list.Items {
			domains = append(domains, strings.Split(x.Spec.Email, "@")[1])
		}
		domains = utils.Unique(domains)

		var suggestions []string
		for _, x := range domains {
			suggestions = append(suggestions, fmt.Sprintf("%s@%s", name, x))
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	return command
}

// Run implements the action
func (o *CreateUserOptions) Run() error {
	user := &orgv1.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "User",
			APIVersion: orgv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: kore.HubNamespace,
		},
		Spec: orgv1.UserSpec{
			Username: o.Name,
			Email:    o.Email,
			Disabled: false,
		},
	}

	return o.WaitForCreation(
		o.Client().
			Resource("user").
			Name(o.Name).
			Payload(user),
		o.NoWait,
	)
}
