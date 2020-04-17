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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/appvia/kore/pkg/cmd"
	authproxy "github.com/appvia/kore/pkg/cmd/auth-proxy"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	o := authproxy.Config{}

	command := &cobra.Command{
		Use:     "auth-proxy",
		Short:   "Kore Authentication Proxy provides a means proxy inbound request to the kube-apiserver",
		Version: version.Version(),

		RunE: func(cmd *cobra.Command, args []string) error {
			var verifiers []authproxy.Verifier

			// @step: create the open id provider
			oidc, err := authproxy.NewOpenIDAuth(
				o.IDPClientID,
				o.IDPServerURL,
				o.UpstreamAuthorizationToken,
				o.IDPUserClaims,
			)
			if err != nil {
				return err
			}
			verifiers = append(verifiers, oidc)

			// @step: create the k8s authorization
			kube, err := authproxy.NewKubeVerifier(o.SigningCA)
			if err != nil {
				return err
			}
			verifiers = append(verifiers, kube)

			svc, err := authproxy.New(log.StandardLogger(), o, verifiers)
			if err != nil {
				return err
			}

			c, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := svc.Run(c); err != nil {
				return err
			}

			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			<-signalChannel

			if err := svc.Stop(); err != nil {
				return err
			}

			return nil
		},
	}

	flags := command.Flags()
	flags.StringVar(&o.Listen, "listen", utils.GetEnvString("LISTEN", ":10443"), "the interface to bind the service to `INTERFACE`")
	flags.StringVar(&o.TLSCert, "tls-cert", utils.GetEnvString("TLS_CERT", ""), "the path to the file containing the certificate pem `PATH`")
	flags.StringVar(&o.TLSKey, "tls-key", utils.GetEnvString("TLS_KEY", ""), "the path to the file containing the private key pem `PATH`")
	flags.StringVar(&o.IDPServerURL, "idp-server-url", utils.GetEnvString("IDP_SERVER_URL", ""), "the open-id server url `URL`")
	flags.StringVar(&o.IDPClientID, "idp-client-id", utils.GetEnvString("IDP_CLIENT_ID", ""), "the identity provider client id used to verify the token `IDP_CLIENT_ID`")
	flags.StringVar(&o.TLSCaAuthority, "ca-authority", utils.GetEnvString("CA_AUTHORITY", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"), "path to kubernetes certificate authority `PATH`")
	flags.BoolVar(&o.EnableProxyProtocol, "enable-proxy-protocol", utils.GetEnvBool("ENABLE_PROXY_PROTOCOL", false), "indicates the proxy should enable proxy protocol support")

	flags.StringSliceVar(&o.IDPUserClaims, "idp-user-claims",
		utils.GetEnvStringSlice("IDP_USER_CLAIMS", []string{"preferred_username", "email", "name"}),
		"an ordered collection of potential token claims to extract the identity `CLAIMS`")

	flags.StringSliceVar(&o.IDPGroupClaims, "idp-group-claims",
		utils.GetEnvStringSlice("IDP_GROUP_CLAIMS", []string{"groups"}),
		"an ordered collection of potential token claims to extract the groups `CLAIMS`")

	flags.StringVar(&o.MetricsListen, "metrics-listen",
		utils.GetEnvString("METRIC_LISTEN", ":8080"), "the interface the prometheus metrics should listen on `INTERFACE`")

	flags.StringVar(&o.UpstreamURL, "upstream-url",
		utils.GetEnvString("UPSTREAM_URL", "https://kubernetes.default.svc.cluster.local"),
		"is the upstream url to forward the requests onto `URL`")

	flags.StringVar(&o.UpstreamAuthorizationToken, "upstream-authentication-token",
		utils.GetEnvString("UPSTREAM_AUTHENTICATION_TOKEN", "/var/run/secrets/kubernetes.io/serviceaccount/token"),
		"the path to the file containing the authentication token for upstream `PATH`")

	flags.Bool("verbose", false, "switches on verbose logging for debugging purposes `BOOL`")
	flags.StringSliceVar(&o.AllowedIPs, "allowed-ips",
		utils.GetEnvStringSlice("ALLOWED_IPS", []string{"0.0.0.0/0"}),
		"traffic will be allowed from the given IP ranges if set. Requires CIDR notation. `CIDR`")

	_ = command.MarkFlagRequired("idp-server-url")
	_ = command.MarkFlagRequired("idp-client-id")
	_ = command.MarkFlagRequired("ca-authority")

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
