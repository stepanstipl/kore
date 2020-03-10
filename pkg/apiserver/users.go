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
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&usersHandler{})
}

type usersHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (u usersHandler) Name() string {
	return "users"
}

// Register is responsible for registering the webserver
func (u *usersHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.WithFields(log.Fields{
		"path": builder.Path("users"),
	}).Info("registering the user webservice with container")

	u.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(builder.Path("users"))

	ws.Route(
		ws.GET("").To(u.findUsers).
			Doc("Returns all the users in the kore").
			Returns(http.StatusOK, "A list of all the users in the kore", orgv1.UserList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{user}").To(u.findUser).
			Doc("Return information related to the specific user in the kore").
			Param(ws.PathParameter("user", "The name of the user you wish to retrieve")).
			Returns(http.StatusOK, "Contains the user definintion from the kore", orgv1.User{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{user}").To(u.updateUser).
			Doc("Used to create or update a user in the kore").
			Param(ws.PathParameter("user", "The name of the user you are updating or creating in the kore")).
			Reads(orgv1.User{}, "The specification for a user in the kore").
			Returns(http.StatusOK, "Contains the user definintion from the kore", orgv1.User{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{user}").To(u.deleteUser).
			Doc("Used to delete a user from the kore").
			Param(ws.PathParameter("user", "The name of the user you are deleting from the kore")).
			Returns(http.StatusOK, "Contains the former user definition from the kore", orgv1.User{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{user}/teams").To(u.findUserTeams).
			Doc("Returns a list of teams the user is a member of").
			Param(ws.PathParameter("user", "The name of the user whos team membership you wish to see")).
			Returns(http.StatusOK, "Response is a team list containing the teams the user is a member of", orgv1.UserList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// findUserTeam returns a list of teams the user is in
func (u usersHandler) findUserTeams(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		user := req.PathParameter("user")

		list, err := u.Users().ListTeams(req.Request.Context(), user)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findUsers returns all the users in the kore
func (u usersHandler) findUsers(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		users, err := u.Users().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, users)
	})
}

// findUser returns a specific user
func (u usersHandler) findUser(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		user, err := u.Users().Get(req.Request.Context(), req.PathParameter("user"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, user)
	})
}

// updateUser is responsible for updating for creating a user in the kore
func (u usersHandler) updateUser(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		user := &orgv1.User{}
		if err := req.ReadEntity(user); err != nil {
			return err
		}

		user, err := u.Users().Update(req.Request.Context(), user)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, user)
	})
}

// deleteUser is responsible for deleting a user from the kore
func (u usersHandler) deleteUser(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		username := req.PathParameter("user")

		user, err := u.Users().Delete(req.Request.Context(), username)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, user)
	})
}
