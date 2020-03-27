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

package options

import (
	"github.com/urfave/cli/v2"
)

// Options returns the command line options
func Options() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "team",
			Aliases: []string{"t"},
			Usage:   "Used to select the team context you are operating in",
		},
		&cli.BoolFlag{
			Name:    "show-flags",
			Usage:   "Used to debugging the flags on the command line `BOOL`",
			Hidden:  true,
			EnvVars: []string{"SHOW_FLAGS"},
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The output format of the resource `FORMAT`",
			Value:   "yaml",
		},
		&cli.BoolFlag{
			Name:  "no-wait",
			Usage: "if we should wait for the resource to provision `BOOL`",
			Value: false,
		},
	}
}
