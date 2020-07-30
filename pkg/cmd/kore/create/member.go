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
	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// CreateMemberOptions is used to provision a team
type CreateMemberOptions struct {
	cmdutils.Factory
	// Invite indicates we generate invitation for the user if required
	Invite bool
	// Team is the team name
	Team string
	// Username is the username to add
	Username string
}

// NewCmdCreateMember returns the create team command
func NewCmdCreateMember(factory cmdutils.Factory) *cobra.Command {
	o := &CreateMemberOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "member",
		Short:   "Adds a user to the team in the kore",
		Example: "kore create member -u <username> [-t team]",
		Run:     cmdutils.DefaultRunFunc(o),
	}
	flags := command.Flags()
	flags.StringVarP(&o.Username, "username", "u", "", "the username of the person being added to the team")
	flags.BoolVarP(&o.Invite, "invite", "i", true, "if the user doesn't exist an invitation url will be automatically generated")

	cmdutils.MustMarkFlagRequired(command, "username")

	cmdutils.MustRegisterFlagCompletionFunc(command, "username", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		suggestions, err := o.Resources().LookResourceNamesWithFilter("user", "", "^"+complete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	return command
}

// Validate is called to validate the options
func (o *CreateMemberOptions) Validate() error {
	if o.Team == "" {
		return errors.ErrTeamMissing
	}

	return nil
}

// Run implements the action
func (o *CreateMemberOptions) Run() error {
	// @step: if we are using invitations we check if the user exists
	if o.Invite {
		found, err := o.ClientWithResource(o.Resources().MustLookup("user")).
			Name(o.Username).
			Exists()
		if err != nil {
			return err
		}
		if !found {
			prompt := promptui.Prompt{
				Label:     "The user does not exist. Do you want to create an invite link",
				IsConfirm: true,
				Default:   "Y",
			}

			if _, err := prompt.Run(); err != nil {
				return nil
			}

			var inviteURL string

			err := o.ClientWithEndpoint("/teams/{team}/invites/generate/{user}").
				Parameters(
					client.PathParameter("team", o.Team),
					client.PathParameter("user", o.Username),
				).
				Result(&inviteURL).
				Get().Error()

			if err != nil {
				return err
			}

			o.Println("Invite URL: %s", inviteURL)

			return nil
		}
	}

	if err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("member")).
		Name(o.Username).
		Update().Error(); err != nil {

		return err
	}
	o.Println("User %q has been added to the team: %s", o.Username, o.Team)

	return nil
}
