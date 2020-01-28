/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"net/http"

	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/hub/authentication"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&whoImpl{})
}

type whoImpl struct {
	hub.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (u whoImpl) Name() string {
	return "users"
}

// Register is responsible for registering the webserver
func (u *whoImpl) Register(i hub.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.WithFields(log.Fields{
		"path": builder.Path("whoami"),
	}).Info("registering the user webservice with container")

	u.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(builder.Path("whoami"))

	ws.Route(
		ws.GET("").To(u.findWho).
			Doc("Returns information about who the user is and what teams they are a member").
			Returns(http.StatusOK, "A list of all the users in the hub", types.WhoAmI{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// findWho checks who you are
func (u *whoImpl) findWho(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		// @step: get the user from context
		user := authentication.MustGetIdentity(req.Request.Context())

		// @step: get the teams
		teams, err := u.Users().ListTeams(req.Request.Context(), user.Username())
		if err != nil {
			return err
		}

		who := &types.WhoAmI{
			Username: user.Username(),
			Email:    user.Email(),
		}
		for _, x := range teams.Items {
			who.Teams = append(who.Teams, x.Name)
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, who)
	})
}
