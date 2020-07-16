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

	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addGCPOrganizationRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/organizations").To(u.findOrganizations).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of gcp organizations").
			Operation("ListGCPOrganizations").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.OrganizationList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/organizations/{name}").To(u.findOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a specific gcp organization").
			Operation("GetGCPOrganization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/organizations/{name}").To(u.updateOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Operation("UpdateGCPOrganization").
			Reads(gcp.Organization{}, "The definition for GCP organization").
			Doc("Is used to provision or update a gcp organization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/organizations/{name}").To(u.deleteOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to delete a managed gcp organization").
			Operation("DeleteGCPOrganization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findOrganizations returns a list of credential claims
func (u teamHandler) findOrganizations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).
			Cloud().GCP().Organizations().
			List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findOrganizations returns a specific credential
func (u teamHandler) findOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).
			Cloud().GCP().Organizations().
			Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateOrganization is used to update an credential claim for a team
func (u teamHandler) updateOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &gcp.Organization{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).
			Cloud().GCP().Organizations().
			Update(req.Request.Context(), claim)

		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteOrganization is used to remove a credential from a team cluster
func (u teamHandler) deleteOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).
			Cloud().GCP().Organizations().
			Delete(req.Request.Context(), name)

		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}
