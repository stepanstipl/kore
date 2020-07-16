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

func (u teamHandler) addGKERoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/gkes").To(u.findGKEs).
			Doc("Returns a list of Google Container Engine clusters which the team has access").
			Operation("ListGKEs").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKEList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkes/{name}").To(u.findGKE).
			Doc("Returns a specific Google Container Engine cluster to which the team has access").
			Operation("GetGKE").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKE{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// findGKEs returns all the clusters under the team
func (u teamHandler) findGKEs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GKE().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findGKE returns a cluster under the team
func (u teamHandler) findGKE(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GKE().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}
