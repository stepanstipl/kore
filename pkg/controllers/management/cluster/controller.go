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

package cluster

import (
	"context"
	"time"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type clusterController struct {
	kore.Interface
	name   string
	logger log.FieldLogger
	// mgr is the manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	ctrl := newClusterController()
	if err := controllers.Register(ctrl); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatalf("failed to register the %s controller", ctrl.Name())
	}
}

func newClusterController() *clusterController {
	name := "cluster"
	return &clusterController{
		name: name,
		logger: log.WithFields(log.Fields{
			"controller": name,
		}),
	}
}

// Name returns the name of the controller
func (a clusterController) Name() string {
	return a.name
}

// Run is called when the controller is started
func (a *clusterController) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	a.Interface = hi

	a.logger.Debug("controller has been started")

	// @step: create the manager for the controller
	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(a))
	if err != nil {
		a.logger.WithError(err).Error("failed to create manager")

		return err
	}

	// @step: set the controller manager
	a.mgr = mgr

	// @step: create the controller
	ctrl, err := controller.New(a.Name(), mgr, controllers.DefaultControllerOptions(a))
	if err != nil {
		a.logger.WithError(err).Error("failed to create the controller")

		return err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(&source.Kind{Type: &clustersv1.Cluster{}},
		&handler.EnqueueRequestForObject{}); err != nil {

		log.WithField("error", err.Error()).Error("failed to create watcher on Cluster resource")

		return err
	}

	go func() {
		a.logger.Info("starting the controller loop")

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

		a.logger.Info("stopping the controller")

		close(a.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (a clusterController) Stop(context.Context) error {
	a.logger.Info("attempting to stop the controller")

	return nil
}
