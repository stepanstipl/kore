/**
 * Copyright (C) 2020 Rohith Jayawardene <info@appvia.io>
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
	"github.com/urfave/cli"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	logger := log.WithFields(log.Fields{
		"config": korectl.HubConfig,
	})

	// @step: load the api config
	config, err := korectl.GetOrCreateClientConfiguration()
	if err != nil {
		logger.WithError(err).Warn("failed to load the kore api configuration")
		logger.Warn("please check the documentation for how to configure the cli")

		os.Exit(1)
	}

	// @step: we need to pull down the swagger and resource cache if required
	if err := korectl.GetCaches(config); err != nil {
		logger.WithError(err).Error("failed to load the cache, try refreshing the cache")

		os.Exit(1)
	}

	app := &cli.App{
		Name:                 "korectl",
		Authors:              version.Authors,
		Author:               version.Prog,
		Email:                version.Email,
		Flags:                options.Options(),
		Usage:                "korectl provides a CLI for the " + version.Prog,
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options %s\n", err)
			return err
		},

		Action: func(ctx *cli.Context) error {
			return nil
		},

		CommandNotFound: func(ctx *cli.Context, name string) {
			fmt.Fprintf(os.Stderr, "[error] command not found %s\n", name)
			os.Exit(1)
		},

		Commands: korectl.GetCommands(config),

		Before: func(ctx *cli.Context) error {
			for _, x := range ctx.Args() {
				for x == "--debug" {
					log.SetLevel(log.DebugLevel)
				}
			}

			command := ctx.Args().Get(0)
			switch {
			case command == "local":
				// no contexts required yet.
			case len(config.Contexts) <= 0:
				log.Warnln("No korectl context configured.")
				os.Exit(0)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
