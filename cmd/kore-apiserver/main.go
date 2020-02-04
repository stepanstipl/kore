/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"os"

	"github.com/appvia/kore/cmd/kore-apiserver/options"
	"github.com/appvia/kore/pkg/cmd"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	app := &cli.App{
		Name:                 "kore-apiserver",
		Authors:              version.Authors,
		Author:               version.Prog,
		Email:                version.Email,
		Flags:                options.Options(),
		Usage:                "Kore API provides the frontend api services for the Appvia Kore",
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options %s\n", err)
			return err
		},

		Action: func(ctx *cli.Context) error {
			return invoke(ctx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
