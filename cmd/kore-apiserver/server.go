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

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/server"
	"github.com/appvia/kore/pkg/services/users"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// invoke is responsible for invoking the api
func invoke(ctx *cli.Context) error {
	// @step: are we enabling verbose logging?
	if ctx.Bool("verbose") {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		log.Debug("enabling verbose logging for debug")
	}
	if ctx.Bool("disable-json-logging") {
		log.SetFormatter(&log.TextFormatter{})
	}

	// @step: construct the server config
	config := server.Config{
		APIServer: apiserver.Config{
			EnableDex:     ctx.Bool("enable-dex"),
			Listen:        ctx.String("listen"),
			MetricsPort:   ctx.Int("metrics-port"),
			SwaggerUIPath: "./swagger-ui",
		},
		Kubernetes: server.KubernetesAPI{
			InCluster:    ctx.Bool("in-cluster"),
			KubeConfig:   ctx.String("kubeconfig"),
			MasterAPIURL: ctx.String("kube-api-server"),
			Token:        ctx.String("kube-api-token"),
		},
		Kore: kore.Config{
			Authenticators:          ctx.StringSlice("kore-authentication-plugin"),
			AdminPass:               ctx.String("admin-pass"),
			AdminToken:              ctx.String("admin-token"),
			ClientID:                ctx.String("client-id"),
			ClientSecret:            ctx.String("client-secret"),
			ClientScopes:            ctx.StringSlice("client-scopes"),
			CertificateAuthority:    ctx.String("certificate-authority"),
			CertificateAuthorityKey: ctx.String("certificate-authority-key"),
			DiscoveryURL:            ctx.String("discovery-url"),
			HMAC:                    ctx.String("kore-hmac"),
			PublicAPIURL:            ctx.String("api-public-url"),
			PublicHubURL:            ctx.String("kore-public-url"),
			UserClaims:              ctx.StringSlice("user-claims"),
			DEX: kore.DEX{
				EnabledDex:    ctx.Bool("enable-dex"),
				PublicURL:     ctx.String("dex-public-url"),
				GRPCServer:    ctx.String("dex-grpc-server"),
				GRPCPort:      ctx.Int("dex-grpc-port"),
				GRPCCaCrt:     ctx.String("dex-grpc-ca-crt"),
				GRPCClientCrt: ctx.String("dex-grpc-client-crt"),
				GRPCClientKey: ctx.String("dex-grpc-client-key"),
			},
		},
		UsersMgr: users.Config{
			EnableLogging: ctx.Bool("enable-user-db-logging"),
			Driver:        ctx.String("users-db-driver"),
			StoreURL:      ctx.String("users-db-url"),
		},
	}

	s, err := server.New(config)
	if err != nil {
		return err
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("attempting to start the kore-apiserver")

	if err := s.Run(c); err != nil {
		return err
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChannel

	// @step: attempt to gracefully stop the api server
	if err := s.Stop(c); err != nil {
		return err
	}

	return nil
}
