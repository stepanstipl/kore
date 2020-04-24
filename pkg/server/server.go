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

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// controller imports
	_ "github.com/appvia/kore/pkg/controllers/register"

	// service provider imports
	_ "github.com/appvia/kore/pkg/serviceproviders"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/services/users"
	"github.com/appvia/kore/pkg/store"
	"github.com/appvia/kore/pkg/utils/crds"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
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
	// rclient is the runtime client
	rclient rc.Client
}

// New is responsible for creating the server container, effectively acting
// as a controller to the other components
func New(config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}

	var kc kubernetes.Interface
	var cc rc.Client

	// register the known types with the schame

	// @step: create the various client
	cfg, err := makeKubernetesConfig(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes config: %s", err)
	}
	kc, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes client: %s", err)
	}

	// @step: ensure we have the kore crds
	crdc, err := crds.NewExtentionsAPIClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create api extensions client: %s", err)
	}
	if err := registerCustomResources(crdc); err != nil {
		return nil, fmt.Errorf("failed to apply the kore crds: %s", err)
	}

	cc, err = rc.New(cfg, rc.Options{Scheme: schema.GetScheme()})
	if err != nil {
		return nil, fmt.Errorf("failed creating runtime client: %s", err)
	}

	// @step: we need to create the data layer
	storecc, err := store.New(kc, cc)
	if err != nil {
		return nil, fmt.Errorf("failed creating store api: %s", err)
	}

	// @step: create the users service
	usermgr, err := users.New(users.Config{
		Driver:        config.UsersMgr.Driver,
		EnableLogging: config.UsersMgr.EnableLogging,
		StoreURL:      config.UsersMgr.StoreURL,
	})
	if err != nil {
		return nil, fmt.Errorf("trying to create the user management service: %s", err)
	}

	// @step: we need to create the kore bridge / business logic
	hubcc, err := kore.New(storecc, usermgr, config.Kore, kore.DefaultServiceProviders)
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
		rclient: cc,
		storecc: storecc,
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
					mgr, err := manager.New(s.cfg, c2.ManagerOptions())
					if err != nil {
						return err
					}

					ctrl, err := controller.New(c2.Name(), mgr, c2.ControllerOptions())
					if err != nil {
						return err
					}

					return c2.RunWithDependencies(ctx, mgr, ctrl, s.hubcc)
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

// Stop is responsible for trying to stop services
func (s serverImpl) Stop(context.Context) error {
	return nil
}
