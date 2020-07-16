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

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addTeamMemberRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/members")).To(u.findTeamMembers).
			Doc("Returns a list of user memberships in the team").
			Operation("ListTeamMembers").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains a collection of team memberships for this team", List{}).
			Returns(http.StatusNotFound, "Team does not exist", nil),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{team}/members/{user}")).To(u.addTeamMember).
			Doc("Used to add a user to the team via membership").
			Operation("AddTeamMember").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are adding to the team")).
			// As there's no body, need to explicitly say we consume any MIME type. Arguably a go-restful bug:
			Consumes(restful.MIME_JSON, "*/*").
			Returns(http.StatusOK, "The user has been successfully added to the team", orgv1.TeamMember{}).
			Returns(http.StatusNotFound, "Team does not exist", nil),
	)

	ws.Route(
		ws.DELETE("/{team}/members/{user}").To(u.removeTeamMember).
			Doc("Used to remove team membership from the team").
			Operation("RemoveTeamMember").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are removing from the team")).
			Returns(http.StatusOK, "The user has been successfully removed from the team", orgv1.TeamMember{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// addTeamMember is used to add a user to a team
func (u teamHandler) addTeamMember(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		user := req.PathParameter("user")

		return u.Teams().Team(team).Members().Add(req.Request.Context(), user)
	})
}

// findTeamMembers returns a list of team membership for this team
func (u teamHandler) findTeamMembers(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		members, err := u.Teams().Team(team).Members().List(req.Request.Context())
		if err != nil {
			return err
		}

		items := makeListWithSize(len(members.Items))
		size := len(members.Items)
		for i := 0; i < size; i++ {
			items.Items[i] = members.Items[i].Spec.Username
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, items)
	})
}

// removeTeamMember is used to remove a user from a team
func (u teamHandler) removeTeamMember(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		user := req.PathParameter("user")

		return u.Teams().Team(team).Members().Delete(req.Request.Context(), user)
	})
}
