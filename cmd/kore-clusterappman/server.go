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
	"os"
	"os/signal"
	"syscall"

	clusterappman "github.com/appvia/kore/pkg/clusterappman"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// invoke is responsible for invoking clusterappman
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
	config := clusterappman.Config{
		Kubernetes: kubernetes.KubernetesAPI{
			InCluster:    ctx.Bool("in-cluster"),
			KubeConfig:   ctx.String("kubeconfig"),
			MasterAPIURL: ctx.String("kube-api-server"),
			Token:        ctx.String("kube-api-token"),
		},
	}

	s, err := clusterappman.New(config)
	if err != nil {
		return err
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("attempting to start the kore cluster manager")

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
