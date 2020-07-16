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

func (u teamHandler) addClusterRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/clusters")).To(u.listClusters).
			Doc("Lists all clusters for a team").
			Operation("ListClusters").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "List of all clusters for a team", clustersv1.ClusterList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/clusters/{name}")).To(u.getCluster).
			Doc("Returns a cluster").
			Operation("GetCluster").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the kubernetes cluster you are acting upon")).
			Returns(http.StatusNotFound, "the cluster with the given name doesn't exist", nil).
			Returns(http.StatusOK, "The requested cluster details", clustersv1.Cluster{}),
	)
	ws.Route(
		withAllErrors(ws.PUT("/{team}/clusters/{name}")).To(u.updateCluster).
			Doc("Creates or updates a cluster").
			Operation("UpdateCluster").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the cluster")).
			Reads(clustersv1.Cluster{}, "The definition for kubernetes cluster").
			Returns(http.StatusOK, "The cluster details", clustersv1.Cluster{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.DELETE("/{team}/clusters/{name}")).To(u.deleteCluster).
			Doc("Deletes a cluster").
			Operation("RemoveCluster").
			Param(ws.PathParameter("name", "Is the name of the cluster")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(params.DeleteCascade()).
			Returns(http.StatusNotFound, "the cluster with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the former cluster definition from the kore", clustersv1.Cluster{}),
	)
}

// listClusters returns all the clusters from a team
func (u teamHandler) listClusters(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Clusters().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getCluster returns a cluster from a team
func (u teamHandler) getCluster(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Clusters().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// updateCluster is responsible for creating or updating a cluster
func (u teamHandler) updateCluster(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		cluster := &clustersv1.Cluster{}
		if err := req.ReadEntity(cluster); err != nil {
			return err
		}

		if err := u.Teams().Team(team).Clusters().Update(req.Request.Context(), cluster); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, cluster)
	})
}

// deleteCluster is responsible for deleting a cluster from a team
func (u teamHandler) deleteCluster(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Clusters().Delete(ctx, name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
