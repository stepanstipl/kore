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

package server

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/appvia/kore/pkg/utils"

	_ "github.com/appvia/kore/pkg/clusterproviders/register"
	_ "github.com/appvia/kore/pkg/controllers/register"
	_ "github.com/appvia/kore/pkg/serviceproviders/register"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/crds"
	korek "github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type serverImpl struct {
	// storecc is the store interface
	storecc store.Store
	// hubcc is the kore interface
	hubcc kore.Interface
	// apicc is the api interface
	apicc apiserver.Interface
	// cfg is the rest.Config for the clients
	cfg *rest.Config
	// client is the runtime client
	client rc.Client
}

// New is responsible for creating the server container, effectively acting
// as a controller to the other components
func New(ctx context.Context, config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// register the known types with the schame

	// @step: create the various client
	cfg, err := korek.MakeKubernetesConfig(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes config: %s", err)
	}
	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes client: %s", err)
	}

	// @step: ensure we have the kore crds
	crdc, err := crds.NewExtentionsAPIClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create api extensions client: %s", err)
	}
	if err := registerCustomResources(ctx, crdc); err != nil {
		return nil, fmt.Errorf("failed to apply the kore crds: %s", err)
	}

	client, err := korek.NewRuntimeClientForAPI(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating runtime client: %s", err)
	}

	// @step: we need to create the data layer
	storecc, err := store.New(kc, client)
	if err != nil {
		return nil, fmt.Errorf("failed creating store api: %s", err)
	}

	// @step: create the persistence service
	persistenceMgr, err := persistence.New(persistence.Config{
		Driver:        config.PersistenceMgr.Driver,
		EnableLogging: config.PersistenceMgr.EnableLogging,
		StoreURL:      config.PersistenceMgr.StoreURL,
	})
	if err != nil {
		return nil, fmt.Errorf("trying to create the user management service: %s", err)
	}

	// @step: we need to create the kore bridge / business logic
	hubcc, err := kore.New(storecc, persistenceMgr, config.Kore)
	if err != nil {
		return nil, fmt.Errorf("trying to create the kore bridge: %s", err)
	}

	if err := makeAuthenticators(hubcc, config); err != nil {
		return nil, err
	}

	// @step: we need to create the apiserver
	apisvr, err := apiserver.New(hubcc, config.APIServer)
	if err != nil {
		return nil, fmt.Errorf("trying to create the apiserver: %s", err)
	}

	return &serverImpl{
		apicc:   apisvr,
		cfg:     cfg,
		hubcc:   hubcc,
		storecc: storecc,
		client:  client,
	}, nil
}

// Run is responsible for starting the services
func (s serverImpl) Run(ctx context.Context) error {

	// @step: we need to start the controllers
	for _, ctrl := range controllers.GetControllers() {
		go func(c controllers.RegisterInterface) {
			log.WithFields(log.Fields{
				"name": c.Name(),
			}).Info("starting the controller")

			err := func() error {
				if c2, ok := c.(controllers.Interface2); ok {
					return s.runController(ctx, c2)
				}

				return c.Run(ctx, s.cfg, s.hubcc)
			}()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"name":  c.Name(),
				}).Fatal("failed to start the controller")
			}
		}(ctrl)
	}

	// @step: start the apiserver - @note this is not being started before
	// the controllers are ready
	if err := s.apicc.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (s serverImpl) runController(ctx context.Context, c controllers.Interface2) error {
	logger := log.WithFields(log.Fields{
		"controller": c.Name(),
	})

	var mopts manager.Options
	if moa, ok := c.(controllers.ManagerOptionsAware); ok {
		mopts = moa.ManagerOptions()
	} else {
		mopts = controllers.DefaultManagerOptions(c)
	}

	mgr, err := manager.New(s.cfg, mopts)
	if err != nil {
		return err
	}

	koreCtx := kore.NewContext(ctx, logger, mgr.GetClient(), s.hubcc)

	var copts controller.Options
	if coa, ok := c.(controllers.ControllerOptionsAware); ok {
		copts = coa.ControllerOptions(koreCtx)
	} else {
		reconciler := reconcile.Func(func(request reconcile.Request) (reconcile.Result, error) {
			logger := koreCtx.Logger().WithFields(log.Fields{
				"name":      request.NamespacedName.Name,
				"namespace": request.NamespacedName.Namespace,
			})
			return c.Reconcile(koreCtx.WithLogger(logger), request)
		})
		copts = controllers.DefaultControllerOptions(reconciler)
	}

	ctrl, err := controller.New(c.Name(), mgr, copts)
	if err != nil {
		return err
	}

	if err := c.Initialize(koreCtx, ctrl); err != nil {
		return err
	}

	go func() {
		logger.Info("starting the controller")

		for {
			if err := mgr.Start(ctx.Done()); err != nil {
				logger.WithError(err).Error("failed to start the controller, retrying")
			}

			if ctx.Err() != nil {
				logger.Info("stopping the controller")
				return
			}

			utils.Sleep(ctx, 5*time.Second)
		}
	}()

	return nil
}

// Stop is responsible for trying to stop services
func (s serverImpl) Stop(context.Context) error {
	return nil
}
