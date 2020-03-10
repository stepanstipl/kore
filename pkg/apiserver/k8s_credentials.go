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

// findKubernetesCredentials returns a list of credential claims
func (u teamHandler) findKubernetesCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).KubernetesCredentials().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findKubernetesCredentials returns a specific credential
func (u teamHandler) findKubernetesCredential(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).KubernetesCredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateKubernetesCredential is used to update an credential claim for a team
func (u teamHandler) updateKubernetesCredential(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &clustersv1.KubernetesCredentials{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).KubernetesCredentials().Update(req.Request.Context(), claim)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteKubernetesCredential is used to remove a credential from a team cluster
func (u teamHandler) deleteKubernetesCredential(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).KubernetesCredentials().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}
