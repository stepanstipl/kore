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

	restful "github.com/emicklei/go-restful"
)

// listKubernetes returns all the Kubernetes clusters under the team
func (u teamHandler) listKubernetes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Kubernetes().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getKubernetes returns a Kubernetes cluster under the team
func (u teamHandler) getKubernetes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Kubernetes().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// deleteKubernetes is responsible for deleting a kubernetes resource
func (u teamHandler) deleteKubernetes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Kubernetes().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateKubernetes is responsible for putting an resource into a team
func (u teamHandler) updateKubernetes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &clustersv1.Kubernetes{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if err := u.Teams().Team(team).Kubernetes().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
