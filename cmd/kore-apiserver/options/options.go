/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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

import "github.com/urfave/cli"

// Options returns the command line options
func Options() []cli.Flag {
	return []cli.Flag{
		//
		// @related to the api server
		//
		cli.StringFlag{
			Name:   "listen",
			Usage:  "the interface to bind the service to `INTERFACE`",
			Value:  ":10080",
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

		//
		// @related to the kore
		//
		cli.StringFlag{
			Name:   "admin-pass",
			Usage:  "the superuser admin password used for first time access",
			EnvVar: "KORE_ADMIN_PASS",
		},
		cli.StringFlag{
			Name:   "admin-token",
			Usage:  "a static admin token which is used to authenticate as admin to the api",
			EnvVar: "KORE_ADMIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "client-id",
			Usage:  "the client id of the openid application `ID`",
			EnvVar: "KORE_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "client-secret",
			Usage:  "the client secret used to setup and oauth2 with dex `SECRET`",
			EnvVar: "KORE_CLIENT_SECRET",
			Value:  "this-should-be-changed",
		},
		cli.StringSliceFlag{
			Name:   "client-scopes",
			Usage:  "additional scopes to add the login request `SCOPE`",
			EnvVar: "KORE_CLIENT_SCOPES",
			Value:  &cli.StringSlice{"profile", "email", "offline"},
		},
		cli.StringFlag{
			Name:   "discovery-url",
			Usage:  "the openid discovery url to use for the openid authenticator `URL`",
			EnvVar: "KORE_DISCOVERY_URL",
		},
		cli.StringFlag{
			Name:   "api-public-url",
			Usage:  "the public url of the api service `URL`",
			EnvVar: "KORE_API_PUBLIC_URL",
		},
		cli.StringFlag{
			Name:   "kore-public-url",
			Usage:  "the public url of the kore service user interface `URL`",
			EnvVar: "KORE_UI_PUBLIC_URL",
		},
		cli.StringFlag{
			Name:     "certificate-authority",
			Usage:    "the path to a file containing a certificate authority `PATH`",
			EnvVar:   "KORE_CERTIFICATE_AUTHORITY",
			Required: true,
		},
		cli.StringFlag{
			Name:     "certificate-authority-key",
			Usage:    "the path to file containing the certificate authority private key  `PATH`",
			EnvVar:   "KORE_CERTIFICATE_AUTHORITY_KEY",
			Required: true,
		},
		cli.StringFlag{
			Name:   "kore-hmac",
			Usage:  "a hmac token used by the kore to sign documents `TOKEN`",
			EnvVar: "KORE_HMAC",
		},
		cli.StringSliceFlag{
			Name:   "kore-authentication-plugin",
			Usage:  "enable one of more authentication plugins for the kore `NAME`",
			EnvVar: "KORE_AUTHENTICATION_PLUGINS",
		},
		cli.StringSliceFlag{
			Name:   "user-claims",
			Usage:  "a list of ordered JWT claims name used to extract the username `NAME`",
			EnvVar: "KORE_USER_CLAIMS",
			Value:  &cli.StringSlice{"preferred_username", "email", "name", "username"},
		},

		//
		// @related to the user management service
		//
		cli.BoolFlag{
			Name:   "enable-user-db-logging",
			Usage:  "enable debug logging on the users and teams database `BOOL`",
			EnvVar: "ENABLE_USER_DB_LOGGING",
		},
		cli.StringFlag{
			Name:   "users-db-driver",
			Usage:  "the database driver which the user managaement service uses `DRIVER`",
			EnvVar: "USERS_DB_DRIVER",
			Value:  "mysql",
		},
		cli.StringFlag{
			Name:   "users-db-url",
			Usage:  "the database dsn used to connect to the users db `DSN`",
			EnvVar: "USERS_DB_URL",
			Value:  "root:pass@tcp(127.0.0.1:3306)/kore?parseTime=true",
		},

		// @related to Dex Identity Provider IDP
		cli.BoolFlag{
			Name:   "enable-dex",
			Usage:  "Indicates if we should enable the dex integration `BOOL`",
			EnvVar: "ENABLE_DEX",
		},
		cli.StringFlag{
			Name:   "dex-public-url",
			Usage:  "the url to the external root of the DEX instance `URL`",
			EnvVar: "DEX_PUBLIC_URL",
			Value:  "http://localhost:5556",
		},
		cli.StringFlag{
			Name:   "dex-grpc-server",
			Usage:  "the remote DEX grpc address `SERVER`",
			EnvVar: "DEX_GRPC_SERVER",
			Value:  "127.0.0.1",
		},
		cli.IntFlag{
			Name:   "dex-grpc-port",
			Usage:  "the remote DEX grpc port `PORT`",
			Value:  5557,
			EnvVar: "DEX_GRPC_PORT",
		},
		cli.StringFlag{
			Name:   "dex-grpc-ca-crt",
			Usage:  "the path to the dex grpc signing ca certificate (optional) `PATH`",
			EnvVar: "DEX_GRPC_CA",
		},
		cli.StringFlag{
			Name:   "dex-grpc-client-crt",
			Usage:  "the path to the dex grpc client certificate (optional) `PATH`",
			EnvVar: "DEX_GRPC_CLIENT_CRT",
		},
		cli.StringFlag{
			Name:   "dex-grpc-client-key",
			Usage:  "the path to the dex grpc client key (optional) `PATH`",
			EnvVar: "DEX_GRPC_CLIENT_KEY",
		},

		//
		// @related to the kubernetes api
		//
		cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "the path to a kubeconfig containing kubernetes config (optional) `PATH`",
		},
		cli.StringFlag{
			Name:   "kube-api-server",
			Usage:  "the url to the kore operations kubernetes api `URL`",
			EnvVar: "KUBE_API_SERVER",
			Value:  "http://127.0.0.1:8080",
		},
		cli.StringFlag{
			Name:   "kube-api-token",
			Usage:  "an optional authorization token for the kube-api `TOKEN`",
			EnvVar: "KUBE_TOKEN",
		},
		cli.BoolFlag{
			Name:   "in-cluster",
			Usage:  "indicates the api is running in the cluster `BOOL`",
			EnvVar: "KUBE_IN_CLUSTER",
		},

		cli.IntFlag{
			Name:   "metrics-port",
			Usage:  "the port the prometheus metrics are served on `PORT`",
			Value:  9090,
			EnvVar: "METRICS_PORT",
		},
		cli.StringFlag{
			Name:   "meta-store-url",
			Usage:  "the url for the meta store i.e. redis://user:pass@hostname `URL`",
			Value:  "",
			EnvVar: "META_STORE_URL",
		},

		// @related to logging
		cli.BoolFlag{
			Name:   "disable-json-logging",
			Usage:  "indicates we should disable json logging `BOOL`",
			EnvVar: "DISABLE_JSON_LOGGING",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "indicates if we should enable verbose logging `BOOL`",
			EnvVar: "VERBOSE",
		},

		// @controller flags
		cli.BoolFlag{
			Name:   "enable-bootstrap-feature",
			Usage:  "Indicates if the bootstrap controller is to be enabled `BOOL`",
			EnvVar: "ENABLE_BOOTSTRAP_FEATURE",
		},
		cli.BoolTFlag{
			Name:   "enable-cluster-deletion-feature",
			Usage:  "Indicates you want the controller delete the cloud resource when deleting the cluster `BOOL`",
			EnvVar: "ENABLE_CLOUD_DELETEION_FEATURE",
		},
	}
}
