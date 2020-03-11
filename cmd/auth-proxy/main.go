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

	"github.com/appvia/kore/cmd/auth-proxy/options"
	"github.com/appvia/kore/pkg/cmd"
	authproxy "github.com/appvia/kore/pkg/cmd/auth-proxy"
	"github.com/appvia/kore/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	app := &cli.App{
		Name: "auth-proxy",
		Authors: []*cli.Author{
			{
				Name:  version.Author,
				Email: version.Email,
			},
		},
		Flags:                options.Options(),
		Usage:                "Kore Authentication Proxy provides a means proxy inbound request to the kube-apiserver",
		Version:              version.Version(),
		EnableBashCompletion: true,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options %s\n", err)
			return err
		},

		Action: func(ctx *cli.Context) error {
			config := authproxy.Config{
				ClientID:                   ctx.String("client-id"),
				DiscoveryURL:               ctx.String("discovery-url"),
				GroupClaims:                ctx.StringSlice("group-claims"),
				Listen:                     ctx.String("listen"),
				MetricsListen:              ctx.String("metrics-listen"),
				SigningCA:                  ctx.String("ca-authority"),
				TLSCert:                    ctx.String("tls-cert"),
				TLSKey:                     ctx.String("tls-key"),
				UserClaims:                 ctx.StringSlice("user-claims"),
				UpstreamURL:                ctx.String("upstream-url"),
				UpstreamAuthorizationToken: ctx.String("upstream-authentication-token"),
			}
			svc, err := authproxy.New(config)
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

	koreCliApp := cmd.NewApp(app)
	if err := koreCliApp.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
