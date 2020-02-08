/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	"context"
	"fmt"

	// controller imports
	_ "github.com/appvia/kore/pkg/controllers/register"

	"github.com/appvia/kore/pkg/apiserver"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/register"
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
	resources, err := register.GetCustomResourceDefinitions()
	if err != nil {
		return nil, fmt.Errorf("failed to decode the kore crds resources: %s", err)
	}
	if err := crds.ApplyCustomResourceDefinitions(crdc, resources); err != nil {
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
	hubcc, err := kore.New(storecc, usermgr, config.Kore)
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
	for _, controller := range controllers.GetControllers() {
		log.WithFields(log.Fields{
			"name": controller.Name(),
		}).Info("starting the controller reconcilation")

		if err := controller.Run(ctx, s.cfg, s.hubcc); err != nil {
			log.WithFields(log.Fields{
				"name":  controller.Name(),
				"error": err.Error(),
			}).Info("failed to start the controller")

			return err
		}
	}

	// @step: start the apiserver
	if err := s.apicc.Run(ctx); err != nil {
		return err
	}

	return nil
}

// Stop is responsible for trying to stop services
func (s serverImpl) Stop(context.Context) error {
	return nil
}
