/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
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

	clusterappman "github.com/appvia/kore/pkg/clusterappman"
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
		Kubernetes: clusterappman.KubernetesAPI{
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
