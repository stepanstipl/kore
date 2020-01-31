/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
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
		// @related to the hub
		//
		cli.StringFlag{
			Name:   "admin-pass",
			Usage:  "the superuser admin password used for first time access",
			EnvVar: "HUB_ADMIN_PASS",
		},
		cli.StringFlag{
			Name:   "admin-token",
			Usage:  "a static admin token which is used to authenticate as admin to the api",
			EnvVar: "HUB_ADMIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "client-id",
			Usage:  "the client id of the openid application `ID`",
			EnvVar: "HUB_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "client-secret",
			Usage:  "the client secret used to setup and oauth2 with dex `SECRET`",
			EnvVar: "HUB_CLIENT_SECRET",
			Value:  "this-should-be-changed",
		},
		cli.StringSliceFlag{
			Name:   "client-scopes",
			Usage:  "additional scopes to add the login request `SCOPE`",
			EnvVar: "HUB_CLIENT_SCOPES",
			Value:  &cli.StringSlice{"profile", "email", "offline"},
		},
		cli.StringFlag{
			Name:   "discovery-url",
			Usage:  "the openid discovery url to use for the openid authenticator `URL`",
			EnvVar: "HUB_DISCOVERY_URL",
		},
		cli.StringFlag{
			Name:   "api-public-url",
			Usage:  "the public url of the api service `URL`",
			EnvVar: "HUB_API_PUBLIC_URL",
		},
		cli.StringFlag{
			Name:   "hub-public-url",
			Usage:  "the public url of the hub service user interface `URL`",
			EnvVar: "HUB_UI_PUBLIC_URL",
		},
		cli.StringFlag{
			Name:     "certificate-authority",
			Usage:    "the path to a file containing a certificate authority `PATH`",
			EnvVar:   "HUB_CERTIFICATE_AUTHORITY",
			Required: true,
		},
		cli.StringFlag{
			Name:     "certificate-authority-key",
			Usage:    "the path to file containing the certificate authority private key  `PATH`",
			EnvVar:   "HUB_CERTIFICATE_AUTHORITY_KEY",
			Required: true,
		},
		cli.StringFlag{
			Name:   "hub-hmac",
			Usage:  "a hmac token used by the hub to sign documents `TOKEN`",
			EnvVar: "HUB_HMAC",
		},
		cli.StringSliceFlag{
			Name:   "hub-authentication-plugin",
			Usage:  "enable one of more authentication plugins for the hub `NAME`",
			EnvVar: "HUB_AUTHENTICATION_PLUGINS",
		},
		cli.StringSliceFlag{
			Name:   "user-claims",
			Usage:  "a list of ordered JWT claims name used to extract the username `NAME`",
			EnvVar: "HUB_USER_CLAIMS",
			Value:  &cli.StringSlice{"preferred_username", "name", "username"},
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
			Value:  "root:pass@tcp(127.0.0.1:3306)/hub?parseTime=true",
		},

		//
		// @related to the user management service
		//
		cli.BoolFlag{
			Name:   "enable-audit-db-logging",
			Usage:  "enables debug logging on the audit and teams database `BOOL`",
			EnvVar: "ENABLE_AUDIT_DB_LOGGING",
		},
		cli.StringFlag{
			Name:   "audit-db-driver",
			Usage:  "the database driver which the user managaement service uses `DRIVER`",
			EnvVar: "AUDIT_DB_DRIVER",
			Value:  "mysql",
		},
		cli.StringFlag{
			Name:   "audit-db-url",
			Usage:  "the database dsn used to connect to the audit db `DSN`",
			EnvVar: "AUDIT_DB_URL",
			Value:  "",
		},

		//
		// @related to Dex Identity Provider IDP
		//
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
			Usage:  "the url to the hub operations kubernetes api `URL`",
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

		//
		// @related to logging
		//
		cli.BoolTFlag{
			Name:   "enable-json-logging",
			Usage:  "indicates we should disable json logging `BOOL`",
			EnvVar: "ENABLE_JSON_LOGGING",
		},
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "indicates if we should enable verbose logging `BOOL`",
			EnvVar: "VERBOSE",
		},
	}
}
