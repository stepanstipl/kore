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

package kore

import (
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

var (
	initLongDescription = `
Allows to you retrieve the resources from the kore api. The command format
is <resource> [name]. When the optional name is not provided we will return
a full listing of all the <resource>s from the API. Examples of resource types
are users, teams, gkes, clusters amongst a few.

You can list all the available resource via $ kore api-resources

Though for a better experience all the resource are autocompletes for you.
Take a look at $ kore completion for details
`
	initExamples = `
# List users:
$ kore get users

#Get information about a specific user:
$ kore get user admin [-o yaml]
`
)

const koreManifestFile = "kore.yml"

// InitOptions the are the options for a get command
type InitOptions struct {
	cmdutil.Factory
}

// NewInitUp creates and returns the up command
func NewInitUp(factory cmdutil.Factory) *cobra.Command {
	o := &UpOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "up",
		Long:    upExamples,
		Example: "up examples",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	//flags := command.Flags()
	//flags.StringVar(&o.Raw, "raw", "", "Raw URI to request from the server")

	return command
}

// Validate is used to validate the options
func (o *InitOptions) Validate() error {
	return nil
}

// Run implements the action
func (o *InitOptions) Run() error {
	path, err := runInit()
	if err != nil {
		return err
	}

	o.Println("Configuration was successfully written to %s", path)

	return nil
}

func runInit() (string, error) {
	return koreManifestFile, nil
}
