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
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// LocalStopOptions is used to stop the demo local environment
type LocalStopOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
}

// NewCmdLocalStop returns the local stop command
func NewCmdLocalStop(factory cmdutil.Factory) *cobra.Command {
	o := &LocalStopOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "stop",
		Short:   "Stops a local demonstration instance of Kore",
		Example: "kore local stop",

		Run: cmdutil.DefaultRunFunc(o),
	}

	return command
}

// Run implements the command action
func (o *LocalStopOptions) Run() error {
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
	cmd, err := getComposeCmd(conf, "down")
	if err != nil {
		return err
	}

	// @step: Execute command
	o.Println("...Stopping Kore.")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		o.Println("%s", stdoutStderr)
		return err
	}

	o.Println("...Kore is now stopped.")
	return nil
}
