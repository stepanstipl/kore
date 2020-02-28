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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli"
)

var longProfileDescription = `

Profiles provide a means to store, configure and switch between multiple
Appvia Kore instances from a single configuration. Alternative, you might 
used profiles to use different profile (i.e. admin / user) to a single 
instance.These are automatically created for you via the $ korectl login 
command or you can manually configure them via the $ korectl profile configure.

Examples: 

$ korectl profile                   # will show this help menu
$ korectl profile show              # will show the profile in use
$ korectl profile list              # will show all the profiles available to you
$ korectl profile use <name>        # switches to another profile
$ korectl profile configure <name>  # allows you to configure a profile 
`

// GetProfilesCommand creates and returns a profiles command
func GetProfilesCommand(config *Config) cli.Command {
	return cli.Command{
		Name:        "profile",
		Usage:       "Used to interact, use, list and show the current profiles available",
		Description: longProfileDescription,

		Subcommands: []cli.Command{
			{
				Name:  "show",
				Usage: "shows the current profile in use",
				Action: func(ctx *cli.Context) error {
					if config.CurrentContext == "" {
						return errors.New("no profiles have create, please you $ korectl login -a <API> or korectl profile configure --help")
					}
					fmt.Println("Profile:  ", config.CurrentContext)
					fmt.Println("Endpoint: ", config.GetCurrentServer().Endpoint)

					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "lists all the profiles which has been configured thus far",
				Action: func(ctx *cli.Context) error {
					// @step: create a tab writer for output
					w := new(tabwriter.Writer)
					w.Init(os.Stdout, 20, 0, 0, ' ', 10)
					defer w.Flush()

					// @step: write out the columns
					_, _ = w.Write([]byte("Name\tServer\n"))

					for _, x := range config.Contexts {
						// @step: ensure the profile has a server and authentication method
						if !config.HasServer(x.Server) || !config.HasAuthInfo(x.AuthInfo) {
							continue
						}
						line := fmt.Sprintf("%s\t%s\n", x.AuthInfo, config.Servers[x.Server].Endpoint)

						_, _ = w.Write([]byte(line))
					}

					return nil
				},
			},
			{
				Name:      "use",
				Usage:     "switches to another profile",
				UsageText: "korectl profile use <name>",
				Action: func(ctx *cli.Context) error {
					if !ctx.Args().Present() {
						return errors.New("you need to specify a profile to use")
					}
					name := ctx.Args().First()

					if !config.HasContext(name) {
						return errors.New("the profile does not exist in your configure")
					}
					config.CurrentContext = name

					if err := config.Update(); err != nil {
						return fmt.Errorf("trying to update your local korectl config: %s", err)
					}

					fmt.Println("Successfully switched the profile to:", name)

					return nil
				},
			},
			{
				Name:      "configure",
				Usage:     "walk through and configure a new profile for you",
				UsageText: "korectl profile configure <name>",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "force the creation of the profile regardless if one exists",
					},
				},
				Action: func(ctx *cli.Context) error {
					// @check they have specified a name
					if !ctx.Args().Present() {
						return errors.New("you need specify a name for the new profile")
					}
					name := ctx.Args().First()

					// @check the profile does not conflict
					if !ctx.Bool("force") && config.HasContext(name) {
						return errors.New("profile name is already taken, please choose another name")
					}

					// @step: ask for the endpoint
					fmt.Printf("Please enter the Kore API: (e.g https://api.domain.com): ")
					endpoint, err := bufio.NewReader(os.Stdin).ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %s", err)
					}
					endpoint = strings.TrimSuffix(endpoint, "\n")

					// @check this is a valid url
					if !IsValidHostname(endpoint) {
						return fmt.Errorf("invalid endpoint: %s", endpoint)
					}

					// @step: create an empty endpoint for then
					config.AddContext(name, &Context{
						Server:   name,
						AuthInfo: name,
					})
					if !config.HasServer(name) {
						config.AddServer(name, &Server{Endpoint: endpoint})
					}
					if !config.HasAuthInfo(name) {
						config.AddAuthInfo(name, &AuthInfo{OIDC: &OIDC{}})
					}
					config.CurrentContext = name

					// @step: attempt to update the configuration
					if err := config.Update(); err != nil {
						return fmt.Errorf("update your local korectl config: %s", err)
					}

					fmt.Println("Successfully configured the profile to: ", name)
					fmt.Println("In order to authenticate please run: $ korectl auth")

					return nil
				},
			},
		},
	}
}
