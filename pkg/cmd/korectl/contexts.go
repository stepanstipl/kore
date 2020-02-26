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
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

func GetContextCommand(config *Config) cli.Command {
	return cli.Command{
		Name:  "context",
		Usage: "Used to manage and interact with the korectl contexts",
		Action: func(_ *cli.Context) error {
			fmt.Printf("Current Context: %s\n", config.CurrentContext)

			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "Used to show a list of contexts available",
				Action: func(ctx *cli.Context) error {
					w := new(tabwriter.Writer)
					w.Init(os.Stdout, 32, 0, 0, ' ', 10)
					defer w.Flush()

					_, _ = w.Write([]byte("Name\tServer\tAuth\n"))
					for _, x := range config.Contexts {
						if !config.HasServer(x.Server) || !config.HasAuthInfo(x.AuthInfo) {
							continue
						}
						_, _ = w.Write([]byte(config.Servers[x.Server].Endpoint + "\t" + x.AuthInfo + "\t\n"))
					}

					return nil
				},
			},
			{
				Name:  "use",
				Usage: "Used to select the current context for the korectl to operate",
				Action: func(ctx *cli.Context) error {
					if !ctx.Args().Present() {
						return errors.New("you need to specify a context to use")
					}
					if !config.HasContext(ctx.Args().First()) {
						return errors.New("the context does not exist")
					}
					config.CurrentContext = ctx.Args().First()

					if err := config.Update(); err != nil {
						return fmt.Errorf("trying to update your locak korectl config: %s", err)
					}
					fmt.Println("successfully switch the context to: ", ctx.Args().First())

					return nil
				},
			},
		},
	}
}
