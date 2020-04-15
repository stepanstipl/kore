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

package delete

import (
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// DeleteAdminOptions the are the options for a delete command
type DeleteAdminOptions struct {
	cmdutils.Factory
	cmdutils.DefaultHandler
	// Username of the user you are removing as an admin
	Username string
}

// NewCmdDeleteAdmin creates and returns the delete admin command
func NewCmdDeleteAdmin(factory cmdutils.Factory) *cobra.Command {
	o := &DeleteAdminOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "admin",
		Short:   "Delete a user from being an admin in kore",
		Example: "kore delete admin --username|-u <username>",
		Run:     cmdutils.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Username, "username", "u", "", "the user you wish to remove from being an admin `USERNAME`")
	cmdutils.MustMarkFlagRequired(command, "username")

	return command
}

// Run implements the action
func (o *DeleteAdminOptions) Run() error {
	err := o.Client().
		Resource("team").
		SubResource("members").
		Name(o.Username).
		Delete().
		Error()
	if err != nil {
		return err
	}
	o.Println("User %q has been successfully removed as an admin")

	return nil
}
