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

	"github.com/urfave/cli/v2"
)

var createLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Example to create a team:
  $ korectl create team a-team
`

// GetCreateCommand creates and returns the create command
func GetCreateCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "create",
		Aliases:     []string{"add"},
		Usage:       "Creates various objects",
		Description: formatLongDescription(createLongDescription),
		ArgsUsage:   "[TYPE] [NAME]",

		Subcommands: []*cli.Command{
			GetCreateTeamCommand(config),
			GetCreateTeamMemberCommand(config),
			GetCreateClusterCommand(config),
			GetCreateNamespaceCommand(config),
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				_ = cli.ShowCommandHelp(ctx, "create")

				return errors.New("[TYPE] [NAME] is required")
			}

			return nil
		},
	}
}
