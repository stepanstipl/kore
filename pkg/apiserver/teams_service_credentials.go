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
	"net/http"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/apiserver/params"
	"github.com/appvia/kore/pkg/kore"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addServiceCredentialRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/servicecredentials")).To(u.listServiceCredentials).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Lists all service credentials for a team").
			Operation("ListServiceCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("cluster", "Is the name of the cluster you are filtering for")).
			Param(ws.QueryParameter("service", "Is the name of the service you are filtering for")).
			Returns(http.StatusOK, "List of all service credentials for a team", servicesv1.ServiceCredentials{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/servicecredentials/{name}")).To(u.getServiceCredentials).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Returns the requsted service credentials").
			Operation("GetServiceCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name of the service credentials")).
			Returns(http.StatusNotFound, "the service credentials with the given name doesn't exist", nil).
			Returns(http.StatusOK, "The requested service crendential details", servicesv1.ServiceCredentials{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{team}/servicecredentials/{name}")).To(u.updateServiceCredentials).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Creates or updates service credentials").
			Operation("UpdateServiceCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the service credentials")).
			Reads(servicesv1.ServiceCredentials{}, "The definition for the service credentials").
			Returns(http.StatusOK, "The service credentail details", servicesv1.ServiceCredentials{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.DELETE("/{team}/servicecredentials/{name}")).To(u.deleteServiceCredentials).
			Filter(filters.FeatureGateFilter(u.Config(), kore.FeatureGateServices)).
			Doc("Deletes the given service credentials").
			Operation("DeleteServiceCredentials").
			Param(ws.PathParameter("name", "Is the name of the service credentials")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(params.DeleteCascade()).
			Returns(http.StatusNotFound, "the service credentials with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the former service credentials definition", servicesv1.ServiceCredentials{}),
	)
}

// listServiceCredentials returns all the service credentials from a team
func (u teamHandler) listServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		var filters []func(servicesv1.ServiceCredentials) bool

		cluster := req.QueryParameter("cluster")
		if cluster != "" {
			filters = append(filters, func(s servicesv1.ServiceCredentials) bool { return s.Spec.Cluster.Name == cluster })
		}

		service := req.QueryParameter("service")
		if service != "" {
			filters = append(filters, func(s servicesv1.ServiceCredentials) bool { return s.Spec.Service.Name == service })
		}

		list, err := u.Teams().Team(team).ServiceCredentials().List(req.Request.Context(), filters...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getServiceCredentials returns service credentials from a team
func (u teamHandler) getServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		list, err := u.Teams().Team(team).ServiceCredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// updateServiceCredentials is responsible for creating or updating service credentials
func (u teamHandler) updateServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		serviceCreds := &servicesv1.ServiceCredentials{}
		if err := req.ReadEntity(serviceCreds); err != nil {
			return err
		}

		if err := u.Teams().Team(team).ServiceCredentials().Update(req.Request.Context(), serviceCreds); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, serviceCreds)
	})
}

// deleteServiceCredentials is responsible for deleting service credentials from a team
func (u teamHandler) deleteServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).ServiceCredentials().Delete(ctx, name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
