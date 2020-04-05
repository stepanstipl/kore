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

	restful "github.com/emicklei/go-restful"
)

// findEKSNodegroups returns all the nodegroups for a EKS cluster for a team
func (u teamHandler) findEKSNodeGroups(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().EKSNodeGroup().List(req.Request.Context())
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findEKS returns a specific nodegroup for a cluster under the team
func (u teamHandler) findEKSNodeGroup(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		ng, err := u.Teams().Team(team).Cloud().EKSNodeGroup().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, ng)
	})
}
