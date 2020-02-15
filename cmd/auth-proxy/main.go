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
	"github.com/urfave/cli"
)

func init() {
	cmd.DefaultLogging()
	log.SetReportCaller(true)
}

func main() {
	app := &cli.App{
		Name:                 "auth-proxy",
		Authors:              version.Authors,
		Author:               version.Prog,
		Email:                version.Email,
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

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
