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

func (u teamHandler) addEKSVPCRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/eksvpcs").To(u.findEKSVPCs).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used to return a list of Amazon EKS VPC which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSVPCList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
		),
	)

	ws.Route(
		ws.GET("/{team}/eksvpcs/{name}").To(u.findEKSVPC).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS VPC you are acting upon")).
			Doc("Is the used to return a EKS VPC which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSVPC{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/eksvpcs/{name}").To(u.updateEKSVPC).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS VPC you are acting upon")).
			Doc("Is used to provision or update a EKS VPC in the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSVPC{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/eksvpcs/{name}").To(u.deleteEKSVPC).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS VPC you are acting upon")).
			Doc("Is used to delete a EKS VPC from the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSVPC{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findEKSVPCs returns all the clusters under the team
func (u teamHandler) findEKSVPCs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKSVPC().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findEKSVPC returns a EKSVPC under the team
func (u teamHandler) findEKSVPC(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKSVPC().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// deleteEKSVPC is responsible for deleting a team resource
func (u teamHandler) deleteEKSVPC(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Cloud().EKSVPC().Get(ctx, name)
		if err != nil {
			return err
		}

		err = u.Teams().Team(team).Cloud().EKSVPC().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateEKSVPC is responsible for putting an resource into a team
func (u teamHandler) updateEKSVPC(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &eks.EKSVPC{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if _, err := u.Teams().Team(team).Cloud().EKSVPC().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
