/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package korectl

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	envsTmpl = path.Join(localCompose, "local.env.tmpl")
)

func prepEnvs(config *Config) error {
	tmpl, err := template.ParseFiles(envsTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create("./demo.env")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config.GetCurrentAuthInfo().OIDC)
}

func startChecks(config *Config) error {
	if !config.HasProfile("local") {
		return errors.New("A 'local' profile has not been found in ~/.korectl/config - try running: korectl local configure.")
	}

	if !config.HasAuthInfo("local") || !config.IsOIDCProviderConfigured("local") {
		return errors.New("No OpenId provider was configured for your 'local' profile in ~/.korectl/config - try running: korectl local configure.")
	}

	return nil
}

func GetLocalRunSubCommands(config *Config) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "start",
			Usage: "Starts a local instance of Kore.",
			Action: func(c *cli.Context) error {
				if err := startChecks(config); err != nil {
					return err
				}

				config.SetCurrentProfile("local")

				if err := prepEnvs(config); err != nil {
					return err
				}

				cmd := exec.Command("docker-compose",
					"--file", "hack/compose/kube.yml",
					"--file", "hack/compose/demo.yml",
					"--file", "hack/compose/operators.yml",
					"up",
					"--force-recreate",
					"-d",
				)
				fmt.Println("...Starting Kore.")

				stdoutStderr, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("%s\n", stdoutStderr)
					return err
				}

				// We pause here to give the services time to initialise
				time.Sleep(time.Second * 35)

				fmt.Printf("...Kore is now started locally and is ready on %s\n", localEndpoint)

				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stops any running local instance of Kore.",
			Action: func(c *cli.Context) error {
				cmd := exec.Command("docker-compose",
					"--file", "hack/compose/kube.yml",
					"--file", "hack/compose/demo.yml",
					"--file", "hack/compose/operators.yml",
					"down",
				)
				fmt.Println("...Stopping Kore.")

				stdoutStderr, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("%s\n", stdoutStderr)
					return err
				}

				fmt.Println("...Kore is now stopped.")

				return nil
			},
		},
	}
}
