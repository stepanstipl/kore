/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

package bootstrap

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/hub"

	clusterv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type bsCtrl struct {
	hub.Interface
	// mgr is the controller manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&bsCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register controller")
	}
}

// Name returns the name of the controller
func (t *bsCtrl) Name() string {
	return "bootstrap-controller"
}

// Run starts the controller
func (t *bsCtrl) Run(ctx context.Context, cfg *rest.Config, i hub.Interface) error {
	t.Interface = i

	logger := log.WithFields(log.Fields{"controller": t.Name()})

	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(t))
	if err != nil {
		logger.WithError(err).Error("failed to create the manager")

		return err
	}
	t.mgr = mgr

	ctrl, err := controller.New(t.Name(), mgr, controllers.DefaultControllerOptions(t))
	if err != nil {
		logger.WithError(err).Error("failed to create the controller")

		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &clusterv1.Kubernetes{}},
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{}); err != nil {

		log.WithError(err).Error("failed to add the controller watcher")
		return err
	}

	go func() {
		logger.Info("starting the controller loop")

		for {
			t.stopCh = make(chan struct{})

			if err := mgr.Start(t.stopCh); err != nil {
				log.WithError(err).Error("failed to start the controller")
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
func (t *bsCtrl) Stop(context.Context) error {
	log.WithField("controller", t.Name()).Info("attempting to stop the controller")

	return nil
}
