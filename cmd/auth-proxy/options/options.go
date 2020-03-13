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
		&cli.StringFlag{
			Name:    "listen",
			Usage:   "the interface to bind the service to `INTERFACE`",
			Value:   ":10443",
			EnvVars: []string{"LISTEN"},
		},
		&cli.StringFlag{
			Name:    "tls-cert",
			Usage:   "the path to the file containing the certificate pem `PATH`",
			EnvVars: []string{"TLS_CERT"},
		},
		&cli.StringFlag{
			Name:    "tls-key",
			Usage:   "the path to the file containing the private key pem `PATH`",
			EnvVars: []string{"TLS_KEY"},
		},
		&cli.StringFlag{
			Name:    "idp-server-url",
			Usage:   "the openid server url `URL`",
			EnvVars: []string{"IDP_SERVER_URL"},
		},
		&cli.StringFlag{
			Name:    "idp-client-id",
			Usage:   "the identity provider client id used to verify the token `IDP_CLIENT_ID`",
			EnvVars: []string{"IDP_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "ca-authority",
			Usage:   "when not using the IDP server url we can use this certificate to verity the token `PATH`",
			EnvVars: []string{"CA_AUTHORITY"},
		},
		&cli.StringFlag{
			Name:    "ca-authority-secret",
			Usage:   "the name of the pre-provision kubernetes secret (namespace/name) which holds the ca `SECRET`",
			EnvVars: []string{"CA_AUTHORITY_SECRET"},
		},
		&cli.StringSliceFlag{
			Name:    "idp-user-claims",
			Usage:   "an ordered collection of potential token claims to extract the identity `CLAIMS`",
			EnvVars: []string{"IDP_USER_CLAIMS"},
			Value:   cli.NewStringSlice("preferred_username", "email", "name"),
		},
		&cli.StringSliceFlag{
			Name:    "idp-group-claims",
			Usage:   "an ordered collection of potential token claims to extract the groups `CLAIMS`",
			EnvVars: []string{"IDP_GROUP_CLAIMS"},
			Value:   cli.NewStringSlice("groups"),
		},
		&cli.StringFlag{
			Name:    "metrics-listen",
			Usage:   "the interface the prometheus metrics should listen on `INTERFACE`",
			EnvVars: []string{"METRICS_LISTEN"},
			Value:   ":8080",
		},
		&cli.StringFlag{
			Name:    "upstream-url",
			Usage:   "is the upstream url to forward the requests onto `URL`",
			EnvVars: []string{"UPSTREAM_URL"},
			Value:   "https://kubernetes.default.svc.cluster.local",
		},
		&cli.StringFlag{
			Name:    "upstream-authentication-token",
			Usage:   "the path to the file containing the authentication token for upstream `PATH`",
			EnvVars: []string{"UPSTREAM_AUTHENTICATION_TOKEN"},
			Value:   "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "switches on verbose logging for debugging purposes `BOOL`",
			EnvVars: []string{"VERBOSE"},
		},
	}
}
