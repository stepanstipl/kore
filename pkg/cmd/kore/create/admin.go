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
	"errors"

	cmdutils "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"

	"github.com/spf13/cobra"
)

// CreateAdminOptions is used to provision a team
type CreateAdminOptions struct {
	cmdutils.Factory
	// Team is the team name
	Team string
	// Username is the username to add
	Username string
}

// NewCmdCreateAdmin returns the create admin command
func NewCmdCreateAdmin(factory cmdutils.Factory) *cobra.Command {
	o := &CreateAdminOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "admin",
		Short:   "Adds to the administator team in kore",
		Example: "kore create admin -u <username> [-t team]",
		Run:     cmdutils.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Username, "username", "u", "", "The username of the person being added to the team")
	command.MarkFlagRequired("username")

	return command
}

// Validate is called to validate the options
func (o *CreateAdminOptions) Validate() error {
	if o.Team == "" {
		return errors.New("no team defined")
	}

	return nil
}

// Run implements the action
func (o *CreateAdminOptions) Run() error {
	if err := o.Client().
		Resource("team").
		Team(kore.HubAdminTeam).
		SubResource("members").
		Name(o.Username).
		Update().Error(); err != nil {

		return err
	}
	o.Println("User %q has been added as a admin", o.Username)

	return nil
}
