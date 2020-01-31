/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

// server is the implementation for the api server
type server struct {
	*Config
	// container is the base container for the services
	container *restful.Container
	// store is the hub bridge interface
	store hub.Interface
}

// New returns a new api server for the hub
func New(hub hub.Interface, config Config) (Interface, error) {
	// @step: verify the configuration
	if err := config.isValid(); err != nil {
		return nil, fmt.Errorf("invalid api config: %s", err)
	}

	// @step: for now we can use the default container
	c := restful.DefaultContainer
	c.Filter(filters.DefaultMetrics.Filter)

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      c,
	}
	c.Filter(cors.Filter)

	authFilter := filters.AuthenticationHandler{
		Realm: hub.Config().DiscoveryURL,
	}

	// @step: register the resource handlers
	for _, x := range GetRegisteredHandlers() {
		ws, err := x.Register(hub, utils.NewPathBuilder(APIVersion))
		if err != nil {
			return nil, err
		}
		if !config.EnableDex && x.Name() == "idp" {
			log.Info("skipping the registation of dex endpoint as feature has been disabled")

			continue
		}
		if x.EnableAuthentication() {
			ws = ws.Filter(authFilter.Filter).Filter(filters.DefaultMembersHandler.Filter)
		}
		if x.EnableLogging() {
			ws = ws.Filter(filters.DefaultLogging.Filter)
		}
		c.Add(ws)
	}

	// @step: register the openapi endpoint service
	c.Add(restfulspec.NewOpenAPIService(restfulspec.Config{
		WebServices:                   c.RegisteredWebServices(),
		APIPath:                       "/swagger.json",
		PostBuildSwaggerObjectHandler: EnrichSwagger,
	}).Filter(filters.SwaggerChecksum.Filter))

	// @step: provide static server of the swagger-ui
	http.Handle("/apidocs/",
		http.StripPrefix("/apidocs/",
			http.FileServer(http.Dir(config.SwaggerUIPath)),
		))

	return &server{Config: &config, container: c, store: hub}, nil
}

// Hub returns the interface to the hub
func (h *server) Hub() hub.Interface {
	return h.store
}

// BaseURI return the base URI
func (h server) BaseURI() string {
	return APIVersion
}

// Run starts the api up
func (h *server) Run(ctx context.Context) error {
	log.WithFields(log.Fields{
		"listen":      h.Listen,
		"tls_enabled": h.UseTLS(),
	}).Info("starting the hub-api")

	// @step: setup the http handler
	s := &http.Server{
		Addr:    h.Listen,
		Handler: h.container,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to start the http server")
		}
	}()

	return nil
}

// Stop indicates to want to stop the api
func (h *server) Stop(ctx context.Context) error {

	return nil
}
