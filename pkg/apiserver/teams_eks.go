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

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addEKSRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/eks").To(u.findEKSs).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used to return a list of Amazon EKS clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/eks/{name}").To(u.findEKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is the used to return a EKS cluster which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKS{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findEKSs returns all the clusters under the team
func (u teamHandler) findEKSs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKS().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findEKS returns a cluster under the team
func (u teamHandler) findEKS(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKS().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}
