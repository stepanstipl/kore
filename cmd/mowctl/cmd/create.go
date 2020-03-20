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
	"github.com/appvia/kore/pkg/cmd/korectl"
	cli "github.com/jawher/mow.cli"
)

func MakeCreateCmd(config *korectl.Config, globals *Globals) func(cmd *cli.Cmd) {
	return func(create *cli.Cmd) {
		create.Command(
			"cluster",
			"Create a kubernetes cluster for a team",
			MakeCreateClusterSubCmd(config, globals),
		)
	}
}
