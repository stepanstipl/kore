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
	"github.com/appvia/kore/pkg/version"
	"github.com/urfave/cli/v2"
)

// Options returns the command line options
func Options() []cli.Flag {
	return []cli.Flag{
		//
		// @related to the api server
		//
		&cli.StringFlag{
			Name:    "listen",
			Usage:   "the interface to bind the service to `INTERFACE`",
			Value:   ":10080",
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

		//
		// @related to the kore
		//
		&cli.StringFlag{
			Name:    "admin-pass",
			Usage:   "the superuser admin password used for first time access",
			EnvVars: []string{"KORE_ADMIN_PASS"},
		},
		&cli.StringFlag{
			Name:    "admin-token",
			Usage:   "a static admin token which is used to authenticate as admin to the api",
			EnvVars: []string{"KORE_ADMIN_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "idp-client-id",
			Usage:   "the client id of the openid application `ID`",
			EnvVars: []string{"KORE_IDP_CLIENT_ID"},
		},
		&cli.StringFlag{
			Name:    "idp-client-secret",
			Usage:   "the client secret used to setup and oauth2 with dex `SECRET`",
			EnvVars: []string{"KORE_IDP_CLIENT_SECRET"},
			Value:   "this-should-be-changed",
		},
		&cli.StringSliceFlag{
			Name:    "idp-client-scopes",
			Usage:   "additional scopes to add the login request `SCOPE`",
			EnvVars: []string{"KORE_IDP_CLIENT_SCOPES"},
			Value:   cli.NewStringSlice("profile", "email", "offline"),
		},
		&cli.StringFlag{
			Name:    "idp-server-url",
			Usage:   "the openid server url to use for the openid authenticator `URL`",
			EnvVars: []string{"KORE_IDP_SERVER_URL"},
		},
		&cli.StringFlag{
			Name:    "local-jwt-public-key",
			Usage:   "the local public key to verify JWTs for localjwt auth plugin",
			EnvVars: []string{"KORE_LOCAL_JWT_PUBLIC_KEY"},
			Value:   "this-should-be-changed",
		},
		&cli.StringFlag{
			Name:    "api-public-url",
			Usage:   "the public url of the api service `URL`",
			EnvVars: []string{"KORE_API_PUBLIC_URL"},
		},
		&cli.StringFlag{
			Name:    "ui-public-url",
			Usage:   "the public url of the kore service user interface `URL`",
			EnvVars: []string{"KORE_UI_PUBLIC_URL"},
		},
		&cli.StringFlag{
			Name:     "certificate-authority",
			Usage:    "the path to a file containing a certificate authority `PATH`",
			EnvVars:  []string{"KORE_CERTIFICATE_AUTHORITY"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "certificate-authority-key",
			Usage:    "the path to file containing the certificate authority private key  `PATH`",
			EnvVars:  []string{"KORE_CERTIFICATE_AUTHORITY_KEY"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "kore-hmac",
			Usage:   "a hmac token used by the kore to sign documents `TOKEN`",
			EnvVars: []string{"KORE_HMAC"},
		},
		&cli.StringSliceFlag{
			Name:    "kore-authentication-plugin",
			Usage:   "enable one of more authentication plugins for the kore `NAME`",
			EnvVars: []string{"KORE_AUTHENTICATION_PLUGINS"},
		},
		&cli.StringSliceFlag{
			Name:    "idp-user-claims",
			Usage:   "a list of ordered JWT claims name used to extract the username `NAME`",
			EnvVars: []string{"KORE_IDP_USER_CLAIMS"},
			Value:   cli.NewStringSlice("preferred_username", "email", "name", "username"),
		},
		&cli.StringFlag{
			Name:    "feature-gates",
			Usage:   "List of feature gates to disable/enable, as key-value pairs, e.g. 'services=true' `GATES`",
			EnvVars: []string{"KORE_FEATURE_GATES"},
		},

		//
		// @related to the user management service
		//
		&cli.BoolFlag{
			Name:    "enable-user-db-logging",
			Usage:   "enable debug logging on the users and teams database `BOOL`",
			EnvVars: []string{"ENABLE_USER_DB_LOGGING"},
		},
		&cli.StringFlag{
			Name:    "users-db-driver",
			Usage:   "the database driver which the user managaement service uses `DRIVER`",
			EnvVars: []string{"USERS_DB_DRIVER"},
			Value:   "mysql",
		},
		&cli.StringFlag{
			Name:    "users-db-url",
			Usage:   "the database dsn used to connect to the users db `DSN`",
			EnvVars: []string{"USERS_DB_URL"},
			Value:   "root:pass@tcp(127.0.0.1:3306)/kore?parseTime=true",
		},

		// @related to Dex Identity Provider IDP
		&cli.BoolFlag{
			Name:    "enable-dex",
			Usage:   "Indicates if we should enable the dex integration `BOOL`",
			EnvVars: []string{"ENABLE_DEX"},
		},
		&cli.StringFlag{
			Name:    "dex-public-url",
			Usage:   "the url to the external root of the DEX instance `URL`",
			EnvVars: []string{"DEX_PUBLIC_URL"},
			Value:   "http://localhost:5556",
		},
		&cli.StringFlag{
			Name:    "dex-grpc-server",
			Usage:   "the remote DEX grpc address `SERVER`",
			EnvVars: []string{"DEX_GRPC_SERVER"},
			Value:   "127.0.0.1",
		},
		&cli.IntFlag{
			Name:    "dex-grpc-port",
			Usage:   "the remote DEX grpc port `PORT`",
			Value:   5557,
			EnvVars: []string{"DEX_GRPC_PORT"},
		},
		&cli.StringFlag{
			Name:    "dex-grpc-ca-crt",
			Usage:   "the path to the dex grpc signing ca certificate (optional) `PATH`",
			EnvVars: []string{"DEX_GRPC_CA"},
		},
		&cli.StringFlag{
			Name:    "dex-grpc-client-crt",
			Usage:   "the path to the dex grpc client certificate (optional) `PATH`",
			EnvVars: []string{"DEX_GRPC_CLIENT_CRT"},
		},
		&cli.StringFlag{
			Name:    "dex-grpc-client-key",
			Usage:   "the path to the dex grpc client key (optional) `PATH`",
			EnvVars: []string{"DEX_GRPC_CLIENT_KEY"},
		},

		//
		// @related to the kubernetes api
		//
		&cli.StringFlag{
			Name:  "kubeconfig",
			Usage: "the path to a kubeconfig containing kubernetes config (optional) `PATH`",
		},
		&cli.StringFlag{
			Name:    "kube-api-server",
			Usage:   "the url to the kore operations kubernetes api `URL`",
			EnvVars: []string{"KUBE_API_SERVER"},
			Value:   "http://127.0.0.1:8080",
		},
		&cli.StringFlag{
			Name:    "kube-api-token",
			Usage:   "an optional authorization token for the kube-api `TOKEN`",
			EnvVars: []string{"KUBE_TOKEN"},
		},
		&cli.BoolFlag{
			Name:    "in-cluster",
			Usage:   "indicates the api is running in the cluster `BOOL`",
			EnvVars: []string{"KUBE_IN_CLUSTER"},
		},

		&cli.IntFlag{
			Name:    "metrics-port",
			Usage:   "the port the prometheus metrics are served on `PORT`",
			Value:   9090,
			EnvVars: []string{"METRICS_PORT"},
		},
		&cli.IntFlag{
			Name:    "profiling-port",
			Usage:   "the port which profiling is service on `PORT`",
			Value:   9091,
			EnvVars: []string{"PROFILING_PORT"},
		},
		&cli.StringFlag{
			Name:    "meta-store-url",
			Usage:   "the url for the meta store i.e. redis://user:pass@hostname `URL`",
			Value:   "",
			EnvVars: []string{"META_STORE_URL"},
		},

		// @related to logging
		&cli.BoolFlag{
			Name:    "disable-json-logging",
			Usage:   "indicates we should disable json logging `BOOL`",
			EnvVars: []string{"DISABLE_JSON_LOGGING"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "indicates if we should enable verbose logging `BOOL`",
			EnvVars: []string{"VERBOSE"},
		},

		// @related to images
		&cli.StringFlag{
			Name:    "auth-proxy-image",
			Usage:   "is the authentication proxy image deployed to the clusters `IMAGE`",
			EnvVars: []string{"AUTH_PROXY_IMAGE"},
			Value:   "quay.io/appvia/auth-proxy:" + version.Release,
		},

		// @controller flags
		&cli.BoolFlag{
			Name:    "enable-cluster-provider-check",
			Value:   true,
			Usage:   "Indicates the kubernetes controller should check the underlying provider status `BOOL`",
			EnvVars: []string{"ENABLE_CLUSTER_PROVIDER_CHECK"},
		},

		&cli.BoolFlag{
			Name:    "enable-profiling",
			Value:   false,
			Usage:   "Indicates we should enable the pprof profile endpoints `BOOL`",
			EnvVars: []string{"ENABLE_PROFILING"},
		},
		&cli.BoolFlag{
			Name:    "enable-metrics",
			Value:   true,
			Usage:   "Indicates we should enable the prometheus metrics `BOOL`",
			EnvVars: []string{"ENABLE_METRICS"},
		},
	}
}
