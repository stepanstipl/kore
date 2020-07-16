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

	aks "github.com/appvia/kore/pkg/apis/aks/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

// addAKSRoutes adds the AKS routes to the web service
func (u teamHandler) addAKSRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/aks").To(u.listAKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used to return a list of AKS clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the definitions of the AKS clusters", aks.AKSList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/aks/{name}").To(u.getAKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the AKS cluster you are acting upon")).
			Doc("Is the used to return the AKS cluster which the team has access").
			Returns(http.StatusOK, "Contains the definition of the AKS cluster", aks.AKS{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// listAKS returns all the AKS clusters under the team
func (u teamHandler) listAKS(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().AKS().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getAKS returns an AKS cluster under the team
func (u teamHandler) getAKS(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().AKS().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}
