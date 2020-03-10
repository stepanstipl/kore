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
