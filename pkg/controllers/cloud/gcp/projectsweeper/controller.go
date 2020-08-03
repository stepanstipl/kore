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

package projectsweeper

import (
	"context"
	"time"

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	"github.com/patrickmn/go-cache"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ controllers.Interface2 = &Controller{}

func init() {
	controllers.Register(NewController())
}

type Controller struct {
	cache *cache.Cache
}

// NewController creates and returns a projectsweeper controller
func NewController() *Controller {
	return &Controller{
		cache: cache.New(30*time.Minute, 1*time.Minute),
	}
}

// Name returns the name of the controller
func (c *Controller) Name() string {
	return "projectsweeper"
}

// ManagerOptions returns the manager options for the controller
func (c *Controller) ManagerOptions() manager.Options {
	options := controllers.DefaultManagerOptions(c)
	options.SyncPeriod = utils.DurationPtr(2 * time.Minute)
	return options
}

// Initialize registers dependencies and sets up watches
func (c *Controller) Initialize(ctx kore.Context, ctrl controller.Controller) error {
	// @step: setup watches for the resources
	if err := ctrl.Watch(
		&source.Kind{Type: &gcp.Project{}},
		&handler.EnqueueRequestForObject{},
	); err != nil {
		ctx.Logger().WithError(err).Error("failed to create watcher on resource")

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
