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

import "github.com/urfave/cli"

// Options returns the command line options
func Options() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "listen",
			Usage:  "the interface to bind the service to `INTERFACE`",
			Value:  ":10443",
			EnvVar: "LISTEN",
		},
		cli.StringFlag{
			Name:   "tls-cert",
			Usage:  "the path to the file containing the certificate pem `PATH`",
			EnvVar: "TLS_CERT",
		},
		cli.StringFlag{
			Name:   "tls-key",
			Usage:  "the path to the file containing the private key pem `PATH`",
			EnvVar: "TLS_KEY",
		},
		cli.StringFlag{
			Name:   "discovery-url",
			Usage:  "the openid discovery url used to pull down idp details `URL`",
			EnvVar: "DISCOVERY_URL",
		},
		cli.StringFlag{
			Name:   "client-id",
			Usage:  "the identity provider client id used to verify the token `CLIENT_ID`",
			EnvVar: "CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "ca-authority",
			Usage:  "when not using the discovery url we can use this certificate to verity the token `PATH`",
			EnvVar: "CA_AUTHORITY",
		},
		cli.StringFlag{
			Name:   "ca-authority-secret",
			Usage:  "the name of the pre-provision kubernetes secret (namespace/name) which holds the ca `SECRET`",
			EnvVar: "CA_AUTHORITY_SECRET",
		},
		cli.StringSliceFlag{
			Name:   "user-claims",
			Usage:  "an ordered collection of potential token claims to extract the identity `CLAIMS`",
			EnvVar: "USER_CLAIMS",
			Value:  &cli.StringSlice{"preferred_username", "email", "name"},
		},
		cli.StringSliceFlag{
			Name:   "group-claims",
			Usage:  "an ordered collection of potential token claims to extract the groups `CLAIMS`",
			EnvVar: "GROUP_CLAIMS",
			Value:  &cli.StringSlice{"groups"},
		},
		cli.StringFlag{
			Name:   "upstream-url",
			Usage:  "is the upstream url to forward the requests onto `URL`",
			EnvVar: "UPSTREAM_URL",
			Value:  "https://kubernetes.default.svc.cluster.local",
		},
		cli.StringFlag{
			Name:   "upstream-authentication-token",
			Usage:  "the path to the file containing the authentication token for upstream `PATH`",
			EnvVar: "UPSTREAM_AUTHENTICATION_TOKEN",
			Value:  "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "switches on verbose logging for debugging purposes `BOOL`",
			EnvVar: "VERBOSE",
		},
	}
}
