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

package local

import (
	"time"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// LocalStartOptions is used to run the demo local environment
type LocalStartOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
}

// NewCmdLocalStart returns the local start command
func NewCmdLocalStart(factory cmdutil.Factory) *cobra.Command {
	o := &LocalStartOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "start",
		Short:   "Starts a local instance of Kore running for demonstration and proof of concept purposes",
		Example: "kore local start",

		Run: cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the command action
func (o *LocalStartOptions) Run() error {
	conf := o.Config()
	// @step: read variables from config file and explode if not present.
	if err := startChecks(conf); err != nil {
		return err
	}
	conf.CurrentProfile = LocalProfileName

	// @step: write out the support files to ~/.korectl/local/
	if err := writeSupportFiles(); err != nil {
		return err
	}

	// @step: prepare command
	cmd, err := getComposeCmd(conf, "up", "--force-recreate", "-d")
	if err != nil {
		return err
	}

	// @step: Execute command
	o.Println("...Starting Kore.")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		o.Println("%s", stdoutStderr)
		return err
	}

	// We pause here to give the services time to initialise
	// @TODO: Perhaps ping for the API and UI for being available instead?
	time.Sleep(time.Second * 35)

	o.Println("...Kore is now started locally")
	o.Println("UI:  http://localhost:3000/")
	o.Println("API: http://localhost:10080/")
	return nil
}
