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

package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/appvia/kore/pkg/controllers/predicates"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var _ controllers.Interface2 = &Controller{}

type Controller struct {
	kore    kore.Interface
	name    string
	logger  log.Ext1FieldLogger
	client  client.Client
	srckind *source.Kind
	kind    string
}

func init() {

	kindsToScan := []*source.Kind{
		{Type: &configv1.Plan{
			TypeMeta: metav1.TypeMeta{
				APIVersion: configv1.GroupVersion.String(),
				Kind:       "Plan",
			},
		}},
		{Type: &clustersv1.Cluster{
			TypeMeta: metav1.TypeMeta{
				APIVersion: clustersv1.GroupVersion.String(),
				Kind:       "Cluster",
			},
		}},
	}

	for _, kind := range kindsToScan {
		ctrl := NewController(log.StandardLogger(), kind)
		if err := controllers.Register(ctrl); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatalf("failed to register the %s controller", ctrl.Name())
		}
	}

}

// NewController creates and returns a new scan controller
func NewController(logger log.Ext1FieldLogger, srckind *source.Kind) *Controller {
	kind := srckind.Type.GetObjectKind().GroupVersionKind().Kind
	name := fmt.Sprintf("security-%s", strings.ToLower(kind))
	return &Controller{
		name: name,
		logger: logger.WithFields(log.Fields{
			"controller": name,
		}),
		srckind: srckind,
		kind:    kind,
	}
}

// Name returns the name of the controller
func (c *Controller) Name() string {
	return c.name
}

func (c *Controller) Logger() log.Ext1FieldLogger {
	return c.logger
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

// Initialize registers dependencies and sets up watches
func (c *Controller) Initialize(ctrl controller.Controller, client client.Client, hi kore.Interface) error {
	c.kore = hi
	c.client = client

	// @step: setup watches for the resources which we support security scanning for
	if err := ctrl.Watch(c.srckind, &handler.EnqueueRequestForObject{}, predicates.SystemResourcePredicate{}); err != nil {
		c.logger.WithError(err).Errorf("failed to create watcher on %s resource", c.kind)
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
