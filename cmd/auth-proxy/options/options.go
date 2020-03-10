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
			Name:    "discovery-url",
			Usage:   "the openid discovery url used to pull down idp details `URL`",
			EnvVars: []string{"DISCOVERY_URL"},
		},
		&cli.StringFlag{
			Name:    "&client-id",
			Usage:   "the identity provider &client id used to verify the token `CLIENT_ID`",
			EnvVars: []string{"CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "ca-authority",
			Usage:   "when not using the discovery url we can use this certificate to verity the token `PATH`",
			EnvVars: []string{"CA_AUTHORITY"},
		},
		&cli.StringFlag{
			Name:    "ca-authority-secret",
			Usage:   "the name of the pre-provision kubernetes secret (namespace/name) which holds the ca `SECRET`",
			EnvVars: []string{"CA_AUTHORITY_SECRET"},
		},
		&cli.StringSliceFlag{
			Name:    "user-claims",
			Usage:   "an ordered collection of potential token claims to extract the identity `CLAIMS`",
			EnvVars: []string{"USER_CLAIMS"},
			Value:   cli.NewStringSlice("preferred_username", "email", "name"),
		},
		&cli.StringSliceFlag{
			Name:    "group-claims",
			Usage:   "an ordered collection of potential token claims to extract the groups `CLAIMS`",
			EnvVars: []string{"GROUP_CLAIMS"},
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
