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

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/appvia/kore/pkg/cmd/korectl"
)

var silentErr = errors.New("silentErr")
var config *korectl.Config
var rootCmd = &cobra.Command{
	Use:   "cobractl",
	Short: "cobractl prototypes building cluster using the Cobra CLI lib.",
	Long: `cobractl prototypes building cluster using the Cobra CLI lib.
This will inform how evolve Kore's korectl. 🐍 CLI power 💪.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize(initConfig)

	// This is required to help with error handling from RunE , https://github.com/spf13/cobra/issues/914#issuecomment-548411337
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return silentErr
	})

	rootCmd.PersistentFlags().StringP(
		"team",
		"t",
		"",
		"Used to select the team context you are operating in",
	)

	rootCmd.PersistentFlags().BoolP(
		"show-flags",
		"",
		viper.GetBool("SHOW_FLAGS"),
		"Used to debugging the flags on the command line `BOOL`",
	)
	rootCmd.PersistentFlags().MarkHidden("show-flags")

	rootCmd.PersistentFlags().StringP(
		"output",
		"o",
		"yaml",
		"The output format of the resource `FORMAT`",
	)
}

func initConfig() {
	config, err := korectl.GetOrCreateClientConfiguration()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read configuration file. Reason: %s\n", err)
		os.Exit(1)
	}

	if err := korectl.GetCaches(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load the cache")
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if err != silentErr {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
