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

func (u teamHandler) addAllocationsRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/allocations").To(u.findAllocations).
			Doc("Used to return a list of all the allocations in the team").
			Operation("ListAllocations").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("assigned", "Retrieves all allocations which have been assigned to you")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.AllocationList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.GET("/{team}/allocations/{name}").To(u.findAllocation).
			Doc("Used to return an allocation within the team").
			Operation("GetAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to return")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/allocations/{name}").To(u.updateAllocation).
			Doc("Used to create/update an allocation within the team.").
			Operation("UpdateAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to update")).
			Reads(configv1.Allocation{}, "The definition of the Allocation").
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.DELETE("/{team}/allocations/{name}").To(u.deleteAllocation).
			Doc("Remove an allocation from a team").
			Operation("RemoveAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to delete")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// findAllocations returns a list of the teams in the allocation
func (u teamHandler) findAllocations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		assigned := req.QueryParameter("assigned")
		kind := req.QueryParameter("kind")
		group := req.QueryParameter("group")

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

		// @step: perform any filtering required
		if kind != "" || group != "" {
			var filtered []configv1.Allocation

			for _, x := range list.Items {
				switch {
				case kind != "" && kind != x.Spec.Resource.Kind:
					continue
				case group != "" && group != x.Spec.Resource.Group:
					continue
				}
				filtered = append(filtered, x)
			}
			list.Items = filtered
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

		if err := u.Teams().Team(team).Allocations().Update(req.Request.Context(), obj, false); err != nil {
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

		obj, err := u.Teams().Team(team).Allocations().Delete(req.Request.Context(), name, false)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, obj)
	})
}
