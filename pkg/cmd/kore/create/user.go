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
	"github.com/appvia/kore/pkg/utils/render"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	userLongDesciption = `
Provides the ability to provision a user in kore.
`

	userExamples = `
# Create the test user
$ kore create user test -e test@appiva.io

# Create a user and provision a local identity for them.
$ kore create user test --password
`
)

// CreateUserOptions is used to provision a team
type CreateUserOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
	// DryRun indicates we only dryrun the resources
	DryRun bool
	// Name is the username to add
	Name string
	// Email is the user email address
	Email string
	// EnableLocal indicates a local identity
	EnableLocal bool
	// LocalPassword creates a user with a local password
	LocalPassword string
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
		Example: userExamples,
		PreRunE: cmdutils.RequireName,
		Run:     cmdutils.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Email, "email", "e", "", "an email address for the user `EMAIL`")
	flags.BoolVar(&o.EnableLocal, "password", false, "used to set a local password on a user `BOOL`")
	flags.BoolVar(&o.DryRun, "dry-run", false, "shows the resource but does not apply or create (defaults: false)")

	cmdutils.MustMarkFlagRequired(command, "email")

	// @step: register the autocompletions
	cmdutils.MustRegisterFlagCompletionFunc(command, "email", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		name := cmd.Flags().Arg(0)
		if name == "" {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		list := &orgv1.UserList{}
		if err := o.ClientWithResource(o.Resources().MustLookup("user")).Result(list).Get().Error(); err != nil {
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
	var err error

	user := makeUserModel(o.Name, o.Email)

	if o.DryRun {
		return render.Render().
			Writer(o.Writer()).
			Format(render.FormatYAML).
			Resource(render.FromStruct(user)).
			Do()
	}

	// @step: check if the user
	found, err := o.ClientWithResource(o.Resources().MustLookup("user")).
		Name(o.Name).
		Exists()
	if err != nil {
		return err
	}
	if !found {
		err = o.WaitForCreation(
			o.ClientWithResource(o.Resources().MustLookup("user")).
				Name(o.Name).
				Payload(user),
			o.NoWait)

		if err != nil {
			return err
		}
	}

	return nil
}

func makeUserModel(username, email string) *orgv1.User {
	return &orgv1.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "User",
			APIVersion: orgv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: kore.HubNamespace,
		},
		Spec: orgv1.UserSpec{
			Username: username,
			Email:    email,
			Disabled: false,
		},
	}
}
