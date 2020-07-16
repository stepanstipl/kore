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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/apiserver/params"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addNamespaceRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/namespaceclaims").To(u.findNamespaces).
			Doc("Used to return all namespaces for the team").
			Operation("ListNamespaces").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former definition from the kore", clustersv1.NamespaceClaimList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/namespaceclaims/{name}").To(u.findNamespace).
			Doc("Used to return the details of a namespace within a team").
			Operation("GetNamespace").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.NamespaceClaim{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/namespaceclaims/{name}").To(u.updateNamespace).
			Doc("Used to create or update the details of a namespace within a team").
			Operation("UpdateNamespace").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Reads(clustersv1.NamespaceClaim{}, "The definition for namespace claim").
			Returns(http.StatusOK, "Contains the definition from the kore", clustersv1.NamespaceClaim{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/namespaceclaims/{name}").To(u.deleteNamespace).
			Doc("Used to remove a namespace from a team").
			Operation("RemoveNamespace").
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(params.DeleteCascade()).
			Returns(http.StatusOK, "Contains the former definition from the kore", clustersv1.NamespaceClaim{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// findNamespaces returns a list of namespace claims
func (u teamHandler) findNamespaces(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).NamespaceClaims().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findNamespace returns a specific namespace
func (u teamHandler) findNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).NamespaceClaims().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateNamespace is used to update an namespace claim for a team
func (u teamHandler) updateNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &clustersv1.NamespaceClaim{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).NamespaceClaims().Update(req.Request.Context(), claim)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteNamespace is used to remove a namespace from a team cluster
func (u teamHandler) deleteNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).NamespaceClaims().Delete(req.Request.Context(), name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}
