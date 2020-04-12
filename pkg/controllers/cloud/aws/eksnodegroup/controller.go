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

package eksnodegroup

import (
	"context"
	"time"

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
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

type ctrl struct {
	kore.Interface
	// mgr is the controller manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&ctrl{}); err != nil {
		log.WithError(err).Fatal("failed to register controller")
	}
}

// Name returns the name of the controller
func (n *ctrl) Name() string {
	return "eksnodegroup"
}

// Run starts the controller
func (n *ctrl) Run(ctx context.Context, cfg *rest.Config, hubi kore.Interface) error {
	logger := log.WithFields(log.Fields{
		"controller": n.Name(),
	})

	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(n))
	if err != nil {
		logger.WithError(err).Error("trying to create the manager")

		return err
	}
	n.mgr = mgr
	n.Interface = hubi

	// @step: create the controller
	ctrl, err := controller.New(n.Name(), mgr, controllers.DefaultControllerOptions(n))
	if err != nil {
		log.WithError(err).Error("failed to create the controller")

		return err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(&source.Kind{Type: &eks.EKSNodeGroup{}},
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{},
	); err != nil {

		log.WithError(err).Error("failed to create watcher on resource")

		return err
	}

	go func() {
		logger.Info("starting the controller loop")

		for {
			n.stopCh = make(chan struct{})

			if err := mgr.Start(n.stopCh); err != nil {
				logger.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()
		logger.Info("stopping the controller")

		close(n.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (n *ctrl) Stop(_ context.Context) error {
	log.WithFields(log.Fields{
		"controller": n.Name(),
	}).Info("attempting to stop the controller")

	return nil
}
