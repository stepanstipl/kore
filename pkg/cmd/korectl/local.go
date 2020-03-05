/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
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

package korectl

import (
	"github.com/urfave/cli"
)

const localEndpoint string = "http://127.0.0.1:10080"
const localManifests string = "./manifests/local"
const localCompose string = "./hack/compose"

func GetLocalCommand(config *Config) cli.Command {
	cmd := cli.Command{
		Name:  "local",
		Usage: "Used to configure and run a local instance of Kore.",
	}
	cmd.Subcommands = append(cmd.Subcommands, GetLocalConfigureSubCommand(config))
	cmd.Subcommands = append(cmd.Subcommands, GetLocalRunSubCommands(config)...)
	cmd.Subcommands = append(cmd.Subcommands, GetLocalLogsSubCommand(config))
	return cmd
}
