/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

package options

import "github.com/urfave/cli/v2"

// Options returns the command line options
func Options() []cli.Flag {
	return []cli.Flag{
		//
		// @related to the kubernetes api
		//
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "the path to a kubeconfig containing kubernetes config (optional) `PATH`",
		},
		&cli.StringFlag{
			Name:    "kube-api-server",
			Usage:   "the url to the hub operations kubernetes api `URL`",
			EnvVars: []string{"KUBE_API_SERVER"},
		},
		&cli.StringFlag{
			Name:    "kube-api-token",
			Usage:   "an optional authorization token for the kube-api `TOKEN`",
			EnvVars: []string{"KUBE_TOKEN"},
		},
		&cli.BoolFlag{
			Name:    "in-cluster",
			Usage:   "indicates the client is running in a cluster `BOOL`",
			EnvVars: []string{"IN_CLUSTER"},
		},

		//
		// @related to logging
		//
		&cli.BoolFlag{
			Name:    "enable-json-logging",
			Value:   true,
			Usage:   "indicates we should disable json logging `BOOL`",
			EnvVars: []string{"ENABLE_JSON_LOGGING"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "indicates if we should enable verbose logging `BOOL`",
			EnvVars: []string{"VERBOSE"},
		},
	}
}
