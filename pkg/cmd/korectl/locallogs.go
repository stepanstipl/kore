/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
				Name:     "follow, f",
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
