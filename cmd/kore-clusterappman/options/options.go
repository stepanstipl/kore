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
