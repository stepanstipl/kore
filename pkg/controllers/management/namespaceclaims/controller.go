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

package namespaceclaims

import (
	"context"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type nsCtrl struct {
	kore.Interface
	// mgr is the manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&nsCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register namespaceclaim controller")
	}
}

// Run is called when the controller is started
func (a *nsCtrl) Run(ctx context.Context, cfg *rest.Config, hi kore.Interface) error {
	a.Interface = hi

	// @step: create the manager for the controller
	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(a))
	if err != nil {
		log.WithError(err).Error("trying to create the manager")

		return err
	}
	a.mgr = mgr
	a.Interface = hi

	// @step: create the controller
	ctrl, err := controller.New(a.Name(), mgr, controllers.DefaultControllerOptions(a))
	if err != nil {
		log.WithError(err).Error("trying to create the controller")

		return err
	}

	// @step: setup watches for the resources
	if err := ctrl.Watch(&source.Kind{Type: &clustersv1.NamespaceClaim{}},
		&handler.EnqueueRequestForObject{},
		&predicate.GenerationChangedPredicate{},
	); err != nil {
		log.WithError(err).Error("failed to create watcher on resource")

		return err
	}

	// @clause: whenever a kubernetes cluster changes we should reconcile the resources
	err = ctrl.Watch(&source.Kind{Type: &clustersv1.Kubernetes{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(o handler.MapObject) []reconcile.Request {
			requests, err := ReconcileNamespaceClaims(ctx, mgr.GetClient(), o.Meta.GetName(), o.Meta.GetNamespace())
			if err != nil {
				log.WithError(err).Error("trying to force reconcilation of namespaceclaims from cluster trigger")

				return []reconcile.Request{}
			}

			return requests
		}),
	})
	if err != nil {
		return err
	}

	// @clause: whenever a team is changed we need queue all namespaceclaims
	// within the team namespace to be reconciled as well
	err = ctrl.Watch(&source.Kind{Type: &orgv1.Team{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(o handler.MapObject) []reconcile.Request {
			requests, err := ReconcileNamespaceClaims(ctx, mgr.GetClient(), o.Meta.GetName(), o.Meta.GetNamespace())
			if err != nil {
				log.WithError(err).Error("trying to force reconcilation of namespaceclaims from team trigger")

				return []reconcile.Request{}
			}

			return requests
		}),
	})
	if err != nil {
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
func (a *nsCtrl) Stop(context.Context) error {
	log.WithFields(log.Fields{
		"controller": a.Name(),
	}).Info("attempting to stop the controller")

	return nil
}

// Name returns the name of the controller
func (a *nsCtrl) Name() string {
	return "namespaceclaims"
}
