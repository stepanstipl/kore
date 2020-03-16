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

// updateTeamSecret is used to add a user to a team
func (u teamHandler) updateTeamSecret(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		secret := &configv1.Secret{}
		if err := req.ReadEntity(secret); err != nil {
			return err
		}
		secret.Name = name

		if err := u.Teams().Team(team).Secrets().Update(req.Request.Context(), secret); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, secret)
	})
}

// findTeamSecrets returns a list of secrets from the team
func (u teamHandler) findTeamSecrets(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		secrets, err := u.Teams().Team(team).Secrets().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, secrets)
	})
}

// findTeamSecret returns a list of secrets from the team
func (u teamHandler) findTeamSecret(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		secret, err := u.Teams().Team(team).Secrets().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, secret)
	})
}

// deleteTeamSecret is used to remove a user from a team
func (u teamHandler) deleteTeamSecret(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		secret, err := u.Teams().Team(team).Secrets().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, secret)
	})
}
