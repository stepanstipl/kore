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
	"fmt"
	"net/http"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/apiserver/params"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addServiceRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/services")).To(u.listServices).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Lists all services for a team").
			Operation("ListServices").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "List of all services for a team", servicesv1.ServiceList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/services/{name}")).To(u.getService).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Returns a service").
			Operation("GetService").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name of the service")).
			Returns(http.StatusNotFound, "the service with the given name doesn't exist", nil).
			Returns(http.StatusOK, "The requested service details", servicesv1.Service{}),
	)
	ws.Route(
		withAllErrors(ws.PUT("/{team}/services/{name}")).To(u.updateService).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Filter(u.readonlyServiceFilter).
			Doc("Creates or updates a service").
			Operation("UpdateService").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the service")).
			Reads(servicesv1.Service{}, "The definition for the service").
			Returns(http.StatusOK, "The service details", servicesv1.Service{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.DELETE("/{team}/services/{name}")).To(u.deleteService).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Filter(u.readonlyServiceFilter).
			Doc("Deletes a service").
			Operation("DeleteService").
			Param(ws.PathParameter("name", "Is the name of the service")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(params.DeleteCascade()).
			Returns(http.StatusNotFound, "the service with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the former service definition from the kore", servicesv1.Service{}),
	)
}

func (u teamHandler) readonlyServiceFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		service, err := u.Teams().Team(team).Services().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if service != nil && service.Annotations[kore.AnnotationReadOnly] == "true" {
			resp.WriteHeader(http.StatusForbidden)
			return nil
		}

		// @step: continue with the chain
		chain.ProcessFilter(req, resp)
		return nil
	})
}

// listServices returns all the services from a team
func (u teamHandler) listServices(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		var list *servicesv1.ServiceList
		var err error

		user := authentication.MustGetIdentity(req.Request.Context())
		if user.IsGlobalAdmin() {
			list, err = u.Teams().Team(team).Services().List(req.Request.Context())
		} else {
			list, err = u.Teams().Team(team).Services().List(req.Request.Context(), func(service servicesv1.Service) bool {
				return service.Annotations[kore.AnnotationSystem] != "true"
			})
		}
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getService returns a service from a team
func (u teamHandler) getService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		service, err := u.Teams().Team(team).Services().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, service)
	})
}

// updateService is responsible for creating or updating a service
func (u teamHandler) updateService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		service := &servicesv1.Service{}
		if err := req.ReadEntity(service); err != nil {
			return err
		}

		if service.Annotations[kore.AnnotationReadOnly] != "" {
			writeError(req, resp, fmt.Errorf("setting %q annotation is not allowed", kore.AnnotationReadOnly), http.StatusForbidden)
			return nil
		}

		if err := u.Teams().Team(team).Services().Update(req.Request.Context(), service); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, service)
	})
}

// deleteService is responsible for deleting a service from a team
func (u teamHandler) deleteService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Services().Delete(ctx, name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
