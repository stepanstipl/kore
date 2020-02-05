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
	"net/http"
	"sync"

	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

var (
	resourceHandlers = &resourceRegistry{}
)

type resourceRegistry struct {
	sync.RWMutex
	// handlers is a collection of handlers
	handlers []Resource
}

// RegisterResourceHandler is called to generate the api from the handler
func RegisterResourceHandlers(hi hub.Interface, builder utils.PathBuilder) ([]*restful.WebService, error) {
	var list []*restful.WebService

	// @step: iterate the handlers
	for _, x := range ResourceHandlers() {
		l := log.WithFields(log.Fields{
			"name": x.Name(),
		})
		l.Info("attempting to generate the rest service to resource name")

		ws := &restful.WebService{}
		ws.Consumes(restful.MIME_JSON)
		ws.Produces(restful.MIME_JSON)
		ws.Path(builder.Path(x.Name()))

		ws.Route(
			ws.GET("").To(func(req *restful.Request, resp *restful.Response) {
				handleErrors(req, resp, func() error {
					list, err := x.List(req.Request.Context())
					if err != nil {
						return err
					}

					return resp.WriteHeaderAndEntity(http.StatusOK, list)
				})
			}).
				Doc("Used to retrieve the a list of the resources from the api").
				Returns(http.StatusOK, "Contains a list of the resources", x.Kind()).
				DefaultReturns("An generic API error containing the cause of the error", Error{}),
		)

		ws.Route(
			ws.GET("/{name}").To(func(req *restful.Request, resp *restful.Response) {
				handleErrors(req, resp, func() error {
					name := req.PathParameter("name")
					resource, err := x.Get(req.Request.Context(), name)
					if err != nil {
						return err
					}

					return resp.WriteHeaderAndEntity(http.StatusOK, resource)
				})
			}).
				Doc("Used to retrieve a specific resource via name from the hub").
				Param(ws.PathParameter("name", "The name of the resource you are retrieve")).
				Returns(http.StatusOK, "Contains the class definintion from the hub", x.Kind()).
				DefaultReturns("An generic API error containing the cause of the error", Error{}),
		)

		ws.Route(
			ws.PUT("/{name}").To(func(req *restful.Request, resp *restful.Response) {
				handleErrors(req, resp, func() error {
					resource := x.Kind()

					if err := req.ReadEntity(resource); err != nil {
						return err
					}

					if err := x.Update(req.Request.Context(), resource); err != nil {
						return err
					}

					return resp.WriteHeaderAndEntity(http.StatusOK, resource)
				})
			}).
				Doc("Used to update the resource in the hub").
				Param(ws.PathParameter("name", "The name of the resource you are updating")).
				Reads(x.Kind(), "The definition for the resource you are updating").
				Returns(http.StatusOK, "Contains the updated resource type", x.Kind()).
				DefaultReturns("An generic API error containing the cause of the error", Error{}),
		)

		ws.Route(
			ws.DELETE("/{name}").To(func(req *restful.Request, resp *restful.Response) {
				handleErrors(req, resp, func() error {
					name := req.PathParameter("name")

					resource, err := x.Get(req.Request.Context(), name)
					if err != nil {
						return err
					}

					if err := x.Delete(req.Request.Context(), name); err != nil {
						return err
					}

					return resp.WriteHeaderAndEntity(http.StatusOK, resource)
				})
			}).
				Doc("Used to delete by name the resource from the hub").
				Param(ws.PathParameter("name", "The name of the resource you are attempting to delete")).
				Returns(http.StatusOK, "Contains the former definition of the resource if any", x.Kind()).
				DefaultReturns("An generic API error containing the cause of the error", Error{}),
		)

		list = append(list, ws)
	}

	return list, nil
}

// Registry is called to register the handler
func (r *resourceRegistry) Register(handler Resource) {
	r.Lock()
	defer r.Unlock()

	r.handlers = append(r.handlers, handler)
}

// Handlers returns the registered handlers
func (r *resourceRegistry) Handlers() []Resource {
	r.RLock()
	defer r.RUnlock()

	return r.handlers
}

// RegisterResource is user to register a resource handler for the storage
func RegisterResource(resource Resource) error {
	resourceHandlers.Register(resource)

	return nil
}

// ResourceHandlers returns the registered handlers
func ResourceHandlers() []Resource {
	return resourceHandlers.Handlers()
}
