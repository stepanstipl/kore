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

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// ShowOptions are the options for the command
type ShowOptions struct {
	cmdutil.Factory
}

// NewCmdProfilesShow creates and returns the profile show command
func NewCmdProfilesShow(factory cmdutil.Factory) *cobra.Command {
	o := &ShowOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "show",
		Aliases: []string{"sh"},
		Short:   "shows the current profile in use",
		Example: "kore profile show",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Validate checks the options
func (o ShowOptions) Validate() error {
	if o.Config().CurrentProfile == "" {
		return fmt.Errorf("no profile selected, use: $ kore profile use <name>")
	}

	return nil
}

// Run implements the action
func (o ShowOptions) Run() error {
	config := o.Config()
	name := config.CurrentProfile

	if !config.HasProfile(config.CurrentProfile) {
		return fmt.Errorf("profile: %s does not exist, $ kore profile use <name>", config.CurrentProfile)
	}

	o.Println("Profile:  %s", name)
	o.Println("Endpoint: %s", config.GetServer(o.Config().CurrentProfile).Endpoint)

	return nil
}
