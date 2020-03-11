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

package korectl

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var editLongDescription = `
The object type accepts both singular and plural nouns (e.g. "user" and "users").

Example to edit a team:
  $ korectl edit team a-team
`

func GetEditCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "edit",
		Usage:       "Modifies various objects",
		Description: formatLongDescription(editLongDescription),
		ArgsUsage:   "[TYPE] [NAME]",

		Subcommands: []*cli.Command{
			GetEditTeamCommand(config),
		},
		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				_ = cli.ShowCommandHelp(ctx.Lineage()[1], "edit")
				fmt.Println()
				return errors.New("[TYPE] [NAME] is required")
			}
			return nil
		},
	}
}
