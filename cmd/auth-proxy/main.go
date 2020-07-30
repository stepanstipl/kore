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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/appvia/kore/pkg/cmd"
	authproxy "github.com/appvia/kore/pkg/cmd/auth-proxy"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers/jwt"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers/openid"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers/tokenreview"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
)

// @TODO the whole codebase needs a cleanup

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

// Config is the conbined configuration
type Config struct {
	// OIDC are the openid options
	OIDC openid.Options
	// LocalJWT are the options for a local jwt verifier
	LocalJWT jwt.Options
	// Server is the main authproxy configuration
	Server authproxy.Config
	// TokenReview are the token review verifier
	TokenReview tokenreview.Options
}

var (
	o = Config{}
)

func main() {
	cmd := &cobra.Command{
		Use:     "auth-proxy",
		Short:   "Authenticates inbound requests to the kube-apiserver",
		Version: version.Version(),

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Validate(); err != nil {
				return err
			}
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				log.SetLevel(log.DebugLevel)
			}

			return Run()
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&o.Server.Listen, "listen", s("LISTEN", ":10443"), "interface to bind the service to")
	flags.StringVar(&o.Server.TLSCert, "tls-cert", s("TLS_CERT", ""), "file containing the certificate pem")
	flags.StringVar(&o.Server.TLSKey, "tls-key", s("TLS_KEY", ""), "file containing the private key pem")
	flags.BoolVar(&o.Server.EnableProxyProtocol, "enable-proxy-protocol", utils.GetEnvBool("ENABLE_PROXY_PROTOCOL", false), "proxy should enable proxy protocol support")
	flags.StringVar(&o.Server.MetricsListen, "metrics-listen", s("METRIC_LISTEN", ":8080"), "the interface prometheus metrics should listen")
	flags.StringSliceVar(&o.Server.AllowedIPs, "allowed-ips", sl("ALLOWED_IPS", []string{"0.0.0.0/0"}), "network cidr allowed access")
	flags.StringVar(&o.Server.UpstreamURL, "upstream-url", s("UPSTREAM_URL", "https://kubernetes.default.svc.cluster.local"), "upstream url to forward the requests")
	flags.StringVar(&o.Server.Token, "upstream-token", s("UPSTREAM_AUTHENTICATION_TOKEN", "/var/run/secrets/kubernetes.io/serviceaccount/token"),
		"containing the authentication token for upstream")
	flags.DurationVar(&o.Server.FlushInterval, "flush-interval", 10*time.Millisecond, "the flush interval used on the revervse proxy")
	flags.StringSliceVar(&o.Server.Verifiers, "verifiers", sl("VERIFIERS", []string{"tokenreview", "localjwt"}), "list of verifiers to enable")
	flags.Bool("verbose", false, "switches on verbose logging for debugging purposes `BOOL`")

	// OpenID options
	flags.StringVar(&o.OIDC.DiscoveryURL, "idp-server-url", s("IDP_SERVER_URL", ""), "the openid server url")
	flags.StringVar(&o.OIDC.ClientID, "idp-client-id", s("IDP_CLIENT_ID", ""), "identity provider client id used to verify the token")
	flags.StringSliceVar(&o.OIDC.UserClaims, "idp-user-claims", sl("IDP_USER_CLAIMS", []string{"preferred_username", "email", "name"}),
		"order list of potential claims to extract the identity `CLAIMS`")

	// LocalJWT options
	flags.StringVar(&o.LocalJWT.ClientID, "jwt-client-id", s("JWT_CLIENT_ID", ""), "client id for the local jwt")
	flags.StringVar(&o.LocalJWT.SignerPath, "jwt-signer-cert", s("JWT_SIGNER_CERT", ""), "path to the certificate used for signing")

	// Tokenreview options
	flags.DurationVar(&o.TokenReview.CacheSuccess, "token-cache-success", 10*time.Minute, "duration to cache successful reviews")
	flags.DurationVar(&o.TokenReview.CacheFailure, "token-cache-failures", 1*time.Minute, "duration to cache failed review")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}

// Validate is responsible for checkinn the inputs
func Validate() error {
	if err := o.OIDC.IsValid(); err != nil {
		return err
	}

	return nil
}

// Run is responsible for implements the action
func Run() error {
	// @step: read in the kubernetes token
	content, err := ioutil.ReadFile(o.Server.Token)
	if err != nil {
		return err
	}
	token := strings.TrimSpace(string(content))

	var verifiers []verifiers.Interface

	if utils.Contains("openid", o.Server.Verifiers) {
		log.Info("initializing the openid verifier")

		o.OIDC.Token = string(token)
		verifier, err := openid.New(o.OIDC)
		if err != nil {
			return err
		}
		verifiers = append(verifiers, verifier)
	}

	if utils.Contains("tokenreview", o.Server.Verifiers) {
		log.Info("initializing the tokenreview verifier")

		config := &rest.Config{
			Host:            o.Server.UpstreamURL,
			BearerToken:     string(token),
			TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		}
		if err != nil {
			return err
		}
		verifier, err := tokenreview.New(config, o.TokenReview)
		if err != nil {
			return err
		}
		verifiers = append(verifiers, verifier)
	}

	if utils.Contains("localjwt", o.Server.Verifiers) {
		log.Info("initializing the localjwt verifier")

		if o.LocalJWT.SignerPath == "" {
			return errors.New("no signing jwt certificate defined")
		}
		signer, err := ioutil.ReadFile(o.LocalJWT.SignerPath)
		if err != nil {
			return err
		}
		o.LocalJWT.Signer = signer
		o.LocalJWT.ImpersonationToken = string(token)

		verifier, err := jwt.New(o.LocalJWT)
		if err != nil {
			return err
		}
		verifiers = append(verifiers, verifier)
	}

	// @step: create the auth proxy service
	svc, err := authproxy.New(&o.Server, verifiers)
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
}

func s(name string, v string) string {
	return utils.GetEnvString(name, v)
}

func sl(name string, v []string) []string {
	return utils.GetEnvStringSlice(name, v)
}
