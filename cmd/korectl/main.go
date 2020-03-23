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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/appvia/kore/cmd/korectl/options"
	"github.com/appvia/kore/pkg/cmd"
	"github.com/appvia/kore/pkg/cmd/korectl"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

// Main is acting as wrapper to the main entrypoint
func Main(args []string, writer, errWriter io.Writer) (int, error) {
	// @step: load the api config
	config, err := korectl.GetOrCreateClientConfiguration()
	if err != nil {
		return 1, fmt.Errorf("failed to read configuration file. reason: %s", err)
	}

	// @step: we need to pull down the swagger and resource cache if required
	if err := korectl.GetCaches(config); err != nil {
		return 1, errors.New("failed to load the cache")
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
		Usage:                "korectl provides a cli for the " + version.Prog,
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			return err
		},

		Commands: korectl.GetCommands(config),

		Action: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				if err := cli.ShowAppHelp(ctx); err != nil {
					return err
				}

				return cli.Exit("", 1)
			}

			return fmt.Errorf(
				"unknown command %q\n\nPlease run `%s --help` to see all available commands.",
				ctx.Args().First(),
				ctx.App.Name,
			)
		},

		Before: func(ctx *cli.Context) error {
			if ctx.Bool("show-flags") {
				fmt.Println("flags:", ctx.Args())
			}

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
			case utils.Contains(command, []string{"profile", "profiles"}):
				return nil
			case command == "login":
				// no contexts required yet.
			case len(config.Profiles) <= 0:
				return fmt.Errorf("no profiles configured.\nPlease check the documentation about how to set up %s.", ctx.App.Name)
			case config.CurrentProfile == "":
				return errors.New("no profile selected.\nPlease use $ korectl profiles --help to select a profile")
			}

			return nil
		},

		Writer:    writer,
		ErrWriter: errWriter,
	}

	koreCliApp := cmd.NewApp(app)
	if err := koreCliApp.Run(args); err != nil {

		switch e := err.(type) {
		case cli.ExitCoder:
			if e.Error() != "" {
				return e.ExitCode(), e
			}

			return e.ExitCode(), nil

		default:
			return 1, err
		}
	}

	return 0, nil
}

func main() {
	exitCode, err := Main(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
