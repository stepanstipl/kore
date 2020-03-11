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
		Name: "korectl",
		Authors: []*cli.Author{
			{
				Name:  version.Author,
				Email: version.Email,
			},
		},
		Flags:                options.Options(),
		Usage:                "korectl provides a CLI for the " + version.Prog,
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			return err
		},

		CommandNotFound: func(ctx *cli.Context, name string) {
			fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", name)
			fmt.Fprintf(os.Stderr, "Please run `%s --help` to see all available commands.\n", ctx.App.Name)
			os.Exit(1)
		},

		Commands: korectl.GetCommands(config),

		Before: func(ctx *cli.Context) error {
			for _, x := range ctx.Args().Slice() {
				for x == "--debug" {
					log.SetLevel(log.DebugLevel)
				}
			}

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

	koreCliApp := cmd.NewApp(app)
	if err := koreCliApp.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		os.Exit(1)
	}
}
