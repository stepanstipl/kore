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

	"github.com/spf13/cobra"
)

var createLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Example to create a team:
  $ korectl create team a-team
`

var createCmd = &cobra.Command{
	Use:     "create [TYPE] [NAME]",
	Aliases: []string{"add"},
	Short:   "Creates various Kore objects",
	Long:    createLongDescription,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("[TYPE] [NAME] is required")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
