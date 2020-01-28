/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package gke

import (
	"context"
	"time"

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/hub"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type gkeCtrl struct {
	hub.Interface
	// mgr is the controller manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&gkeCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register controller")
	}
}

// Name returns the name of the controller
func (t *gkeCtrl) Name() string {
	return "gke.compute.hub.appvia.io"
}

// Run starts the controller
func (t *gkeCtrl) Run(ctx context.Context, cfg *rest.Config, hubi hub.Interface) error {
	logger := log.WithFields(log.Fields{
		"controller": t.Name(),
	})

	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(t))
	if err != nil {
		logger.WithError(err).Error("trying to create the manager")

		return err
	}
	t.mgr = mgr
	t.Interface = hubi

	// @step: create the controller for gke
	if _, err = controllers.NewController(
		"gke.hub.appvia.io", t.mgr,
		&source.Kind{Type: &gke.GKE{}},
		&controllers.ReconcileHandler{
			HandlerFunc: t.Reconcile,
		},
	); err != nil {
		log.WithError(err).Error("trying to create the managed cluster roles controller")

		return err
	}

	if _, err = controllers.NewController(
		"gke-credentials.hub.appvia.io", t.mgr,
		&source.Kind{Type: &gke.GKECredentials{}},
		&controllers.ReconcileHandler{
			HandlerFunc: t.ReconcileCredentials,
		},
	); err != nil {
		log.WithError(err).Error("trying to create the managed cluster roles controller")

		return err
	}

	go func() {
		logger.Info("starting the controller loop")

		for {
			t.stopCh = make(chan struct{})

			if err := mgr.Start(t.stopCh); err != nil {
				logger.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()
		logger.Info("stopping the controller")

		close(t.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (t *gkeCtrl) Stop(_ context.Context) error {
	log.WithFields(log.Fields{
		"controller": t.Name(),
	}).Info("attempting to stop the controller")

	return nil
}
