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
	"regexp"
	"strings"
	"syscall"

	"github.com/appvia/kore/pkg/utils/kubernetes"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// invoke is responsible for invoking the api
func invoke(ctx *cli.Context) error {
	// @step: are we enabling verbose logging?
	if ctx.Bool("disable-json-logging") {
		log.SetFormatter(&log.TextFormatter{})
	}
	if ctx.Bool("verbose") {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		log.Debug("enabling verbose logging for debug")
	}

	featuresGates, err := mergeFeatureGates(ctx.String("feature-gates"))
	if err != nil {
		return err
	}

	config := server.Config{
		APIServer: apiserver.Config{
			EnableDex:       ctx.Bool("enable-dex"),
			EnableMetrics:   ctx.Bool("enable-metrics"),
			EnableProfiling: ctx.Bool("enable-profiling"),
			Listen:          ctx.String("listen"),
			MetricsPort:     ctx.Int("metrics-port"),
			ProfilingPort:   ctx.Int("profiling-port"),
			SwaggerUIPath:   "./swagger-ui",
		},
		Kubernetes: kubernetes.KubernetesAPI{
			InCluster:    ctx.Bool("in-cluster"),
			KubeConfig:   ctx.String("kubeconfig"),
			MasterAPIURL: ctx.String("kube-api-server"),
			Token:        ctx.String("kube-api-token"),
		},
		Kore: kore.Config{
			AdminPass:                  ctx.String("admin-pass"),
			AdminToken:                 ctx.String("admin-token"),
			AuthProxyImage:             ctx.String("auth-proxy-image"),
			Authenticators:             ctx.StringSlice("kore-authentication-plugin"),
			CertificateAuthority:       ctx.String("certificate-authority"),
			CertificateAuthorityKey:    ctx.String("certificate-authority-key"),
			EnableClusterProviderCheck: ctx.Bool("enable-cluster-provider-check"),
			FeatureGates:               featuresGates,
			HMAC:                       ctx.String("kore-hmac"),
			IDPClientID:                ctx.String("idp-client-id"),
			IDPClientScopes:            ctx.StringSlice("idp-client-scopes"),
			IDPClientSecret:            ctx.String("idp-client-secret"),
			IDPServerURL:               ctx.String("idp-server-url"),
			IDPUserClaims:              ctx.StringSlice("idp-user-claims"),
			LocalJWTPublicKey:          ctx.String("local-jwt-public-key"),
			PublicAPIURL:               ctx.String("api-public-url"),
			PublicHubURL:               strings.TrimRight(ctx.String("ui-public-url"), "/"),
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
		PersistenceMgr: persistence.Config{
			EnableLogging: ctx.Bool("enable-user-db-logging"),
			Driver:        ctx.String("users-db-driver"),
			StoreURL:      ctx.String("users-db-url"),
		},
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := server.New(c, config)
	if err != nil {
		return err
	}

	log.WithField("feature-gates", config.Kore.FeatureGates).Info("starting kore-apiserver")

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

var featureGateRegexp = regexp.MustCompile(`^\s*(.*?)\s*=\s*(true|false)\s*$`)

func mergeFeatureGates(values string) (map[string]bool, error) {
	featureGates := kore.DefaultFeatureGates()

	if strings.TrimSpace(values) == "" {
		return featureGates, nil
	}

	for _, s := range strings.Split(values, ",") {
		parts := featureGateRegexp.FindStringSubmatch(s)
		if len(parts) == 0 {
			return nil, fmt.Errorf("--feature-gates is invalid, it must be in 'service1=true,service2=false' format")
		}
		gate := parts[1]
		enabled := parts[2] == "true"

		if _, exists := featureGates[gate]; !exists {
			return nil, fmt.Errorf("feature gate %q does not exist", gate)
		}

		featureGates[gate] = enabled
	}

	return featureGates, nil
}
