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

package main

import (
	"fmt"
	"os"

	"github.com/appvia/kore/cmd/korectl/options"
	"github.com/appvia/kore/pkg/cmd"
	"github.com/appvia/kore/pkg/cmd/korectl"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	// @step: load the api config
	config, err := korectl.GetOrCreateClientConfiguration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read configuration file. Reason: %s\n", err)
		os.Exit(1)
	}

	// @step: we need to pull down the swagger and resource cache if required
	if err := korectl.GetCaches(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load the cache")
		os.Exit(1)
	}

	app := &cli.App{
		Name:                 "korectl",
		Flags:                options.Options(),
		Usage:                "korectl provides a CLI for the " + version.Prog,
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			return err
		},

		CommandNotFound: func(ctx *cli.Context, name string) {
			fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", name)
			fmt.Fprintf(os.Stderr, "Please run `%s help` to see all available commands.\n", ctx.App.Name)
			os.Exit(1)
		},

		Commands: korectl.GetCommands(config),

		Before: func(ctx *cli.Context) error {
			command := ctx.Args().Get(0)
			if command == "" || ctx.App.Command(command) == nil { // We don't have a valid command
				return nil
			}

			switch {
			case command == "local", command == "help", command == "autocomplete":
				// no contexts required yet.
			case command == "profile" && ctx.Args().Get(1) == "configure":
				// no contexts required yet.
			case command == "login":
				// no contexts required yet.
			case len(config.Profiles) <= 0:
				fmt.Fprintf(os.Stderr, "Error: no %s profiles configured.\n", ctx.App.Name)
				fmt.Fprintf(os.Stderr, "Please check the documentation about how to set up %s.\n", ctx.App.Name)
				os.Exit(1)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		os.Exit(1)
	}
}
