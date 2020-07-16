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

func (u teamHandler) addGCPProjectClaimRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/projectclaims").To(u.findProjectClaims).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaimList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/projectclaims/{name}").To(u.findProjectClaim).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findProjectClaims returns a list of credential claims
func (u teamHandler) findProjectClaims(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GCP().ProjectClaims().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findProjectClaims returns a specific credential
func (u teamHandler) findProjectClaim(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).Cloud().GCP().ProjectClaims().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}
