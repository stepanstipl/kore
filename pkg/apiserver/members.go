/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
