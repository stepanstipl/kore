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

	"github.com/appvia/kore/pkg/utils"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/crds"
	korek "github.com/appvia/kore/pkg/utils/kubernetes"

	// controller imports
	_ "github.com/appvia/kore/pkg/controllers/register"

	// service provider imports
	_ "github.com/appvia/kore/pkg/serviceproviders/register"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type serverImpl struct {
	// config is the server configuration
	config Config
	// hubcc is the kore interface
	hubcc kore.Interface
	// apicc is the api interface
	apiServer apiserver.Interface
	// cfg is the rest.Config for the clients
	restConfig *rest.Config
	// client is the runtime client
	client rc.Client
}

// New is responsible for creating the server container, effectively acting
// as a controller to the other components
func New(ctx context.Context, config Config, persistenceMgr persistence.Interface) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	// @step: create the various client
	restConfig, err := korek.MakeKubernetesConfig(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes config: %s", err)
	}
	kc, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes client: %s", err)
	}

	// @step: ensure we have the kore crds
	crdc, err := crds.NewExtentionsAPIClient(restConfig)
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

	storecc, err := store.New(kc, client)
	if err != nil {
		return nil, fmt.Errorf("failed creating store api: %s", err)
	}

	if persistenceMgr == nil {
		// @step: create the persistence service
		persistenceMgr, err = persistence.New(persistence.Config{
			Driver:        config.PersistenceMgr.Driver,
			EnableLogging: config.PersistenceMgr.EnableLogging,
			StoreURL:      config.PersistenceMgr.StoreURL,
		})
		if err != nil {
			return nil, fmt.Errorf("trying to create the user management service: %s", err)
		}
	}

	// @step: we need to create the kore bridge / business logic
	hubcc, err := kore.New(storecc, persistenceMgr, config.Kore)
	if err != nil {
		return nil, fmt.Errorf("trying to create the kore bridge: %s", err)
	}

	if err := makeAuthenticators(hubcc, config); err != nil {
		return nil, err
	}

	var apiServer apiserver.Interface
	if config.APIServer.Enabled {
		if apiServer, err = apiserver.New(hubcc, config.APIServer); err != nil {
			return nil, fmt.Errorf("trying to create the apiserver: %s", err)
		}
	}

	return &serverImpl{
		config:     config,
		apiServer:  apiServer,
		restConfig: restConfig,
		hubcc:      hubcc,
		client:     client,
	}, nil
}

func (s serverImpl) enabledControllers() []controllers.RegisterInterface {
	var res []controllers.RegisterInterface
	for _, c := range controllers.GetControllers() {
		if utils.Contains(c.Name(), s.config.Kore.Controllers) {
			res = append(res, c)
		}
	}
	return res
}

// Run is responsible for starting the services
func (s serverImpl) Run(ctx context.Context) error {
	// @step: we need to start the controllers
	for _, ctrl := range s.enabledControllers() {
		go func(c controllers.RegisterInterface) {
			log.WithFields(log.Fields{
				"name": c.Name(),
			}).Info("starting the controller")

			err := func() error {
				if c2, ok := c.(controllers.Interface2); ok {
					mgr, err := manager.New(s.restConfig, c2.ManagerOptions())
					if err != nil {
						return err
					}

					ctrl, err := controller.New(c2.Name(), mgr, c2.ControllerOptions())
					if err != nil {
						return err
					}

					return c2.RunWithDependencies(ctx, mgr, ctrl, s.hubcc)
				}

				return c.Run(ctx, s.restConfig, s.hubcc)
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
	if s.config.APIServer.Enabled {
		if err := s.apiServer.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop is responsible for trying to stop services
func (s serverImpl) Stop(ctx context.Context) error {
	if s.config.APIServer.Enabled {
		if err := s.apiServer.Stop(ctx); err != nil {
			log.WithError(err).Error("stopping the API server")
		}
	}

	for _, ctrl := range s.enabledControllers() {
		if err := ctrl.Stop(ctx); err != nil {
			log.WithError(err).Errorf("stopping %s controller", ctrl.Name())
		}
	}

	return nil
}
