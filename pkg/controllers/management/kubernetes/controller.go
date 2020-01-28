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

package kubernetes

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/hub"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type k8sCtrl struct {
	hub.Interface
	// mgr is the manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&k8sCtrl{}); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to register the kubernetes controller")
	}
}

// Name returns the name of the controller
func (a k8sCtrl) Name() string {
	return "kubernetes-cluster"
}

// Run is called when the controller is started
func (a *k8sCtrl) Run(ctx context.Context, cfg *rest.Config, hi hub.Interface) error {
	a.Interface = hi

	logger := log.WithFields(log.Fields{
		"controller": a.Name(),
	})
	logger.Debug("controller has been started")

	// @step: create the manager for the controller
	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(a))
	if err != nil {
		logger.WithError(err).Error("trying to create manager")

		return err
	}

	// @step: set the controller manager
	a.mgr = mgr

	// @step: create the controller
	ctrl, err := controller.New(a.Name(), mgr, controllers.DefaultControllerOptions(a))
	if err != nil {
		logger.WithError(err).Error("trying to create the controller")

		return err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(&source.Kind{Type: &clustersv1.Kubernetes{}},
		&handler.EnqueueRequestForObject{}); err != nil {

		log.WithField("error", err.Error()).Error("failed to create watcher on resource")

		return err
	}

	// @step: watch for any changes to the teams, find all clusters they have an
	err = ctrl.Watch(&source.Kind{Type: &orgv1.Team{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(o handler.MapObject) []reconcile.Request {
			requests, err := ReconcileClusterRequests(ctx, mgr.GetClient(), o.Meta.GetName())
			if err != nil {
				log.Error(err, "failed to force reconcilation of clusters from trigger")

				return []reconcile.Request{}
			}

			return requests
		})})
	if err != nil {
		return err
	}

	go func() {
		logger.Info("starting the controller loop")

		for {
			a.stopCh = make(chan struct{})

			if err := mgr.Start(a.stopCh); err != nil {
				log.WithField(
					"error", err.Error(),
				).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()

		logger.Info("stopping the controller")

		close(a.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (a k8sCtrl) Stop(context.Context) error {
	log.WithFields(log.Fields{
		"controller": a.Name(),
	}).Info("attempting to stop the controller")

	return nil
}
