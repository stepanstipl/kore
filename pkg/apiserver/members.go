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
