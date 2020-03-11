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

package korectl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

const koreApiServer = "/kore-apiserver"

var dcomposeArgs = []string{
	"--file", "hack/compose/kube.yml",
	"--file", "hack/compose/demo.yml",
	"--file", "hack/compose/operators.yml",
}

func GetLocalLogsSubCommand(_ *Config) *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "View logs from local Kore.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "follow",
				Aliases:  []string{"f"},
				Usage:    "Follow log output.",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("...Retrieving Kore logs.")

			running, err := isKoreStarted()
			if err != nil {
				return err
			}

			if !running {
				fmt.Println("...No logs are available (Kore is not running).")
				return nil
			}

			return runLogs(c)
		},
	}
}

func isKoreStarted() (bool, error) {
	ps := append(dcomposeArgs, "ps")
	cmd := exec.Command("docker-compose", ps...)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", stdoutStderr)
		return false, err
	}

	return strings.Contains(string(stdoutStderr), koreApiServer), nil
}

func runLogs(c *cli.Context) error {
	logs := append(dcomposeArgs, "logs")

	if c.Bool("follow") {
		logs = append(logs, "--follow")
	}

	cmd := exec.Command("docker-compose", logs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
