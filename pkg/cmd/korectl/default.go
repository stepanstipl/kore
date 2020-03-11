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

import "github.com/urfave/cli/v2"

// DefaultsOptions are options for all commands
var DefaultOptions = []cli.Flag{
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "The output format of the resource `FORMAT`",
		Value:   "yaml",
	},
	&cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"D"},
		Usage:   "Indicates for verbose logging `BOOL`",
	},
}
