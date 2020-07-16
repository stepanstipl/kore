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

func (u teamHandler) addEKSCredentialsRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/ekscredentials").To(u.listEKSCredentials).
			Operation("ListEKSCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Amazon EKS credentials which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentialsList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/ekscredentials/{name}").To(u.getEKSCredentials).
			Operation("GetEKSCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS Credentials you are acting upon")).
			Doc("Is the used tor return a list of EKS Credentials which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/ekscredentials/{name}").To(u.updateEKSCredentials).
			Operation("UpdateEKSCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS credentials you are acting upon")).
			Reads(eks.EKSCredentials{}, "The definition for EKS Credentials").
			Doc("Is used to provision or update a EKS credentials in the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/ekscredentials/{name}").To(u.deleteEKSCredentials).
			Operation("DeleteEKSCredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS credentials you are acting upon")).
			Doc("Is used to delete a EKS credentials from the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// listEKSCredentials returns all the credentials under the team
func (u teamHandler) listEKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKSCredentials().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getEKSCredentials returns credentials under the team
func (u teamHandler) getEKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKSCredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// deleteEKSCredentials is responsible for deleting a team resource
func (u teamHandler) deleteEKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Cloud().EKSCredentials().Get(ctx, name)
		if err != nil {
			return err
		}

		err = u.Teams().Team(team).Cloud().EKSCredentials().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateEKSCredentials is responsible for putting an resource into a team
func (u teamHandler) updateEKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &eks.EKSCredentials{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if _, err := u.Teams().Team(team).Cloud().EKSCredentials().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
