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

	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&whoImpl{})
}

type whoImpl struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (u whoImpl) Name() string {
	return "users"
}

// Register is responsible for registering the webserver
func (u *whoImpl) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Path("whoami")

	log.WithFields(log.Fields{
		"path": path,
	}).Info("registering the user webservice with container")

	u.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path)

	ws.Route(
		withStandardErrors(ws.GET("")).To(u.findWho).
			Doc("Returns information about who the user is and what teams they are a member").
			Operation("WhoAmI").
			Returns(http.StatusOK, "A list of all the users in the kore", types.WhoAmI{}),
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
			AuthMethod: user.AuthMethod(),
			Username:   user.Username(),
			Email:      user.Email(),
		}
		for _, x := range teams.Items {
			who.Teams = append(who.Teams, x.Name)
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, who)
	})
}
