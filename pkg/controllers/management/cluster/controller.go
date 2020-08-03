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

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ controllers.Interface2 = &Controller{}

type Controller struct {
	kore.Interface
	name   string
	logger log.Ext1FieldLogger
	client client.Client
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
func NewController(logger log.Ext1FieldLogger) *Controller {
	name := "cluster"
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

func (c *Controller) Logger() log.Ext1FieldLogger {
	return c.logger
}

// ManagerOptions returns the manager options for the controller
func (c *Controller) ManagerOptions() manager.Options {
	return controllers.DefaultManagerOptions(c)
}

// ControllerOptions returns the controllers options
func (c *Controller) ControllerOptions() controller.Options {
	return controllers.DefaultControllerOptions(c)
}

// Initialize registers dependencies and sets up watches
func (c *Controller) Initialize(ctrl controller.Controller, client client.Client, hi kore.Interface) error {
	c.client = client
	c.Interface = hi

	if err := ctrl.Watch(
		&source.Kind{Type: &clustersv1.Cluster{}},
		&handler.EnqueueRequestForObject{},
		predicate.GenerationChangedPredicate{},
	); err != nil {
		c.logger.WithError(err).Error("failed to create watcher on Cluster resource")
		return err
	}

	return nil
}

func (c *Controller) Run(context.Context, *rest.Config, kore.Interface) error {
	panic("deprecated")
}

func (c *Controller) Stop(context.Context) error {
	panic("deprecated")
}
