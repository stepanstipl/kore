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

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/validation"
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
	// store is the kore bridge interface
	store kore.Interface
}

// withStandardErrors adds the standard internal server error (500) result to the route.
func withStandardErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{})
}

// withValidationErrors adds the standard bad request (400) validation error result to the route.
func withValidationErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusBadRequest, "Validation error of supplied parameters/body", validation.Error{})
}

// withAuthErrors adds the standard unauthenticated (401) and forbidden (403) results to the route.
func withAuthErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return rb.
		Returns(http.StatusUnauthorized, "If not authenticated", nil).
		Returns(http.StatusForbidden, "If authenticated but not authorized", nil)
}

// withAllErrors is a shorthand to add all standard, validation, and auth results to the route.
func withAllErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return withValidationErrors(withAuthErrors(withStandardErrors(rb)))
}

// withAllNonValidationErrors is a shorthand to add all standard and auth results to the route but not validation.
func withAllNonValidationErrors(rb *restful.RouteBuilder) *restful.RouteBuilder {
	return withAuthErrors(withStandardErrors(rb))
}

// New returns a new api server for the kore
func New(kore kore.Interface, config Config) (Interface, error) {
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
		Realm: kore.Config().IDPServerURL,
	}

	// @step: register the resource handlers
	for _, x := range GetRegisteredHandlers() {
		ws, err := x.Register(kore, utils.NewPathBuilder(APIVersion))
		if err != nil {
			return nil, err
		}
		if x.EnableAuthentication() {
			ws = ws.Filter(authFilter.Filter).Filter(filters.DefaultMembersHandler.Filter)
		}
		if x.EnableAudit() {
			// Register the auditing filter on a per-route basis so we can audit the
			// operation name.
			routes := ws.Routes()
			for idx := range routes {
				routes[idx].Filters = append(routes[idx].Filters,
					filters.NewAuditingFilter(
						kore.Audit,
						APIVersion,
						ws.RootPath(),
						routes[idx].Operation))
			}
		}
		if x.EnableLogging() {
			ws = ws.Filter(filters.DefaultLogging.Filter)
		}
		if !x.Enabled() {
			ws = ws.Filter(filters.DefaultNotImplementedHandler.Filter)
		}
		c.Add(ws)
	}

	// @step: register the openapi endpoint service
	c.Add(restfulspec.NewOpenAPIService(restfulspec.Config{
		APIPath:                       "/swagger.json",
		PostBuildSwaggerObjectHandler: EnrichSwagger,
		WebServices:                   c.RegisteredWebServices(),
	}).Filter(filters.SwaggerChecksum.Filter))

	// @step: provide static server of the swagger-ui
	http.Handle("/apidocs/",
		http.StripPrefix("/apidocs/",
			http.FileServer(http.Dir(config.SwaggerUIPath)),
		))

	return &server{Config: &config, container: c, store: kore}, nil
}

// Kore returns the interface to the kore
func (h *server) Kore() kore.Interface {
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
	}).Info("starting the kore-apiserver")

	// @step: setup the http handler
	s := &http.Server{Addr: h.Listen, Handler: h.container}

	go func() {
		var err error

		switch h.UseTLS() {
		case true:
			err = s.ListenAndServeTLS(h.TLSCert, h.TLSKey)
		default:
			err = s.ListenAndServe()
		}
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to start the http server")
		}
	}()

	return nil
}

// Stop indicates to want to stop the api
func (h *server) Stop(ctx context.Context) error {
	log.Info("attempting to stop the kore-apiserver")

	return nil
}
