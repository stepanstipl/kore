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

package features

import (
	"context"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ controllers.Interface2 = &Controller{}

// Controller is the features controller
type Controller struct {
	kore   kore.Interface
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

// NewController creates and returns a new features controller
func NewController(logger log.FieldLogger) *Controller {
	name := "korefeatures"
	return &Controller{
		name: name,
		logger: logger.WithFields(log.Fields{
			"controller": name,
		}),
	}
}

// Name returns the name of the controller
func (c *Controller) Name() string {
	return c.name
}

// ManagerOptions are the manager options
func (c *Controller) ManagerOptions() manager.Options {
	options := controllers.DefaultManagerOptions(c)
	options.SyncPeriod = utils.DurationPtr(3 * time.Hour)

	return options
}

// ControllerOptions are the controller options
func (c *Controller) ControllerOptions() controller.Options {
	return controllers.DefaultControllerOptions(c)
}

// RunWithDependencies provisions the controller
func (c *Controller) RunWithDependencies(ctx context.Context, mgr manager.Manager, ctrl controller.Controller, ki kore.Interface) error {
	c.kore = ki
	c.mgr = mgr
	c.ctrl = ctrl

	c.logger.Debug("controller has been started")

	if err := c.ctrl.Watch(
		&source.Kind{Type: &configv1.KoreFeature{}},
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{},
	); err != nil {
		c.logger.WithError(err).Error("failed to create watcher on resource")
		return err
	}

	var stopCh chan struct{}

	// Start a loop to repeatedly start the manager until the context is cancelled.
	go func() {
		c.logger.Info("starting the controller loop")
		for {
			stopCh = make(chan struct{})
			if err := c.mgr.Start(stopCh); err != nil {
				c.logger.WithError(err).Error("failed to start the controller")
			}
			if ctx.Err() != nil {
				// Context was cancelled
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Monitor for the context completing and, on completion, send the stop signal to the
	// manager by closing the current stopCh.
	go func() {
		<-ctx.Done()
		c.logger.Info("stopping the controller")
		if stopCh != nil {
			close(stopCh)
		}
	}()

	return nil
}

// Run is called when the controller is started
func (c *Controller) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	panic("this controller implements controllers.Interface2 and only RunWithDependencies should be called")
}

// Stop is responsible for calling a halt on the controller
func (c *Controller) Stop(context.Context) error {
	c.logger.Info("attempting to stop the controller")

	return nil
}