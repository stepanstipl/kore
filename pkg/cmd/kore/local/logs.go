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
	"os"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

// LocalLogsOptions is used to run the demo local environment
type LocalLogsOptions struct {
	cmdutil.Factory
	cmdutil.DefaultHandler
	// Follow indicates we should follow the logs rather than just outputting them
	Follow bool
}

// NewCmdLocalLogs returns the local logs command
func NewCmdLocalLogs(factory cmdutil.Factory) *cobra.Command {
	o := &LocalLogsOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "logs",
		Short:   "Shows the current logs ",
		Example: "kore local logs",

		Run: cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.BoolVarP(&o.Follow, "follow", "f", false, "set to follow the log `BOOL`")

	return command
}

// Run implements the command action
func (o *LocalLogsOptions) Run() error {
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

	running, err := isKoreStarted(conf)
	if err != nil {
		return err
	}
	if !running {
		o.Println("...No logs are available (Kore is not running).")
		return nil
	}

	logsArgs := []string{"logs"}
	if o.Follow {
		logsArgs = append(logsArgs, "--follow")
	}
	logsCmd, err := getComposeCmd(conf, logsArgs...)
	if err != nil {
		return err
	}

	logsCmd.Stdout = os.Stdout
	logsCmd.Stderr = os.Stderr

	return logsCmd.Run()
}
