/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package deployments

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	appsv1 "github.com/appvia/kore/pkg/apis/apps/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type dpCtrl struct {
	kore.Interface
	// mgr is the manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&dpCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register the allocations controller")
	}
}

// Name returns the name of the controller
func (a dpCtrl) Name() string {
	return "deployments.kore.appvia.io"
}

// Run is called when the controller is started
func (a *dpCtrl) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	a.Interface = hi

	// @step: create the manager for the controller
	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(a))
	if err != nil {
		log.WithError(err).Error("failed to create the manager")

		return err
	}

	// @step: set the controller manager
	a.mgr = mgr

	// @step: create an application controller
	if _, err := controllers.NewController(finalizerName, a.mgr,
		&source.Kind{Type: &appsv1.AppDeployment{}},
		&controllers.ReconcileHandler{
			HandlerFunc: a.Reconcile,
		},
	); err != nil {
		log.WithError(err).Error("trying to create global application controller")

		return err
	}

	go func() {
		log.Info("starting the controller loop")

		for {
			a.stopCh = make(chan struct{})

			if err := mgr.Start(a.stopCh); err != nil {
				log.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()
		log.WithFields(log.Fields{
			"controller": a.Name(),
		}).Info("stopping the controller")

		close(a.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (a dpCtrl) Stop(context.Context) error {
	log.WithFields(log.Fields{
		"controller": a.Name(),
	}).Info("attempting to stop the controller")

	return nil
}
