/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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

		// @NOTE: this will do for now
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
