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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	restful "github.com/emicklei/go-restful"
)

// findAllocations returns a list of the teams in the allocation
func (u teamHandler) findAllocations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		assigned := req.QueryParameter("assigned")

		if assigned == "false" {
			list, err := u.Teams().Team(team).Allocations().List(req.Request.Context())
			if err != nil {
				return err
			}

			return resp.WriteHeaderAndEntity(http.StatusOK, list)
		}

		list, err := u.Teams().Team(team).Allocations().ListAllocationsAssigned(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findAllocation returns an allocation in the team
func (u teamHandler) findAllocation(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")
		assigned := req.QueryParameter("assigned")

		if assigned == "false" {
			obj, err := u.Teams().Team(team).Allocations().Get(req.Request.Context(), name)
			if err != nil {
				return err
			}

			return resp.WriteHeaderAndEntity(http.StatusOK, obj)
		}

		obj, err := u.Teams().Team(team).Allocations().GetAssigned(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, obj)
	})
}

// updateAllocation is responsible for updating the allocations
func (u teamHandler) updateAllocation(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		obj := &configv1.Allocation{}
		if err := req.ReadEntity(obj); err != nil {
			return err
		}
		obj.Name = name

		if err := u.Teams().Team(team).Allocations().Update(req.Request.Context(), obj); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, obj)
	})
}

// deleteAllocation removes any allocations from the team
func (u teamHandler) deleteAllocation(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		obj, err := u.Teams().Team(team).Allocations().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, obj)
	})
}
