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

	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addGKECredentialsRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/gkecredentials").To(u.findGKECredientalss).
			Doc("Returns a list of GKE Credentials to which the team has access").
			Operation("ListGKECredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentialsList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkecredentials/{name}").To(u.findGKECredientals).
			Doc("Returns a specific GKE Credential to which the team has access").
			Operation("GetGKECredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/gkecredentials/{name}").To(u.updateGKECredientals).
			Doc("Creates or updates a specific GKE Credential to which the team has access").
			Operation("UpdateGKECredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Reads(gke.GKECredentials{}, "The definition for GKE Credentials").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/gkecredentials/{name}").To(u.deleteGKECredientals).
			Doc("Deletes a specific GKE Credential from the team").
			Operation("DeleteGKECredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// findGKECredientalss returns all the clusters under the team
func (u teamHandler) findGKECredientalss(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GKECredentials().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findGKECredientals returns a cluster under the team
func (u teamHandler) findGKECredientals(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GKECredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// deleteGKECredientals is responsible for deleting a team resource
func (u teamHandler) deleteGKECredientals(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Cloud().GKECredentials().Get(ctx, name)
		if err != nil {
			return err
		}

		err = u.Teams().Team(team).Cloud().GKECredentials().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateGKECredientals is responsible for putting an resource into a team
func (u teamHandler) updateGKECredientals(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &gke.GKECredentials{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if _, err := u.Teams().Team(team).Cloud().GKECredentials().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
