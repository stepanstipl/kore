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

package namespaces

import (
	"context"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ controllers.Interface2 = &Controller{}

// Controller implements the reconcilations logic
type Controller struct {
	kore.Interface
	name   string
	logger log.FieldLogger
	mgr    manager.Manager
	ctrl   controller.Controller
}

func init() {
	ctrl := NewController(log.StandardLogger())
	if err := controllers.Register(ctrl); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatalf("failed to register the %s controller", ctrl.Name())
	}
}

// NewController creates and returns a clusters controller
func NewController(logger log.FieldLogger) *Controller {
	name := "namespaces-imports"
	return &Controller{
		name: name,
		logger: logger.WithFields(log.Fields{
			"controller": name,
		}),
	}
}

// Name returns the name of the controller
func (a *Controller) Name() string {
	return a.name
}

// ManagerOptions returns the manager options for the controller
func (a *Controller) ManagerOptions() manager.Options {
	resync := 5 * time.Minute
	options := controllers.DefaultManagerOptions(a)
	options.RetryPeriod = &resync

	return options
}

// ControllerOptions returns the controllers options
func (a *Controller) ControllerOptions() controller.Options {
	return controllers.DefaultControllerOptions(a)
}

// RunWithDependencies starts the controller up
func (a *Controller) RunWithDependencies(ctx context.Context, mgr manager.Manager, ctrl controller.Controller, hi kore.Interface) error {
	a.mgr = mgr
	a.ctrl = ctrl
	a.Interface = hi

	a.logger.Debug("controller has been started")

	// @step: setup watches for the resources
	if err := a.ctrl.Watch(
		&source.Kind{Type: &clustersv1.Cluster{}},
		&handler.EnqueueRequestForObject{},
		predicate.GenerationChangedPredicate{},
	); err != nil {

		a.logger.WithError(err).Error("failed to create watcher on Cluster resource")

		return err
	}

	var stopCh chan struct{}

	go func() {
		a.logger.Info("starting the controller loop")

		for {
			stopCh = make(chan struct{})

			if err := a.mgr.Start(stopCh); err != nil {
				a.logger.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()

		a.logger.Info("stopping the controller")

		if stopCh != nil {
			close(stopCh)
		}
	}()

	return nil
}

// Run is called when the controller is started
func (a *Controller) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	panic("this controller implements controllers.Interface2 and only RunWithDependencies should be called")
}

// Stop is responsible for calling a halt on the controller
func (a *Controller) Stop(context.Context) error {
	a.logger.Info("attempting to stop the controller")

	return nil
}
