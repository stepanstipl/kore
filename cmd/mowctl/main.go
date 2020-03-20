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

package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/appvia/kore/pkg/cmd/korectl"
)

type CmdGlobals struct {
	Team, Output string
	ShowFlags    bool
}

var cmdGlobals CmdGlobals

func main() {
	config, err := korectl.GetOrCreateClientConfiguration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read configuration file. Reason: %s\n", err)
		os.Exit(1)
	}

	if err := korectl.GetCaches(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load the cache")
		os.Exit(1)
	}

	app := cli.App("mowctl", "mowctl prototypes building a cluster using the mow.cli lib.")
	app.LongDesc = `mowctl prototypes building cluster using the mow.cli lib.
This will inform how we evolve Kore's korectl - CLI power 💪.`

	app.StringPtr(&cmdGlobals.Team, cli.StringOpt{
		Name:      "t team",
		Desc:      "Used to select the team context you are operating in",
		EnvVar:    "",
		Value:     "",
		HideValue: false,
		SetByUser: nil,
	})
	app.StringPtr(&cmdGlobals.Team, cli.StringOpt{
		Name:      "o output",
		Desc:      "The output format of the resource `FORMAT`",
		HideValue: false,
		SetByUser: nil,
	})
	app.BoolPtr(&cmdGlobals.ShowFlags, cli.BoolOpt{
		Name:      "show-flags",
		Desc:      "Used to debugging the flags on the command line `BOOL`",
		EnvVar:    "SHOW_FLAGS",
		Value:     false,
		HideValue: true,
		SetByUser: nil,
	})

	if err := app.Run(os.Args); err != nil {
		cli.Exit(1)
	}
}
