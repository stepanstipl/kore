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

func (u teamHandler) addTeamSecretRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/secrets").To(u.findTeamSecrets).
			Operation("ListTeamSecrets").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to return all the secrets within the team").
			Returns(http.StatusOK, "Contains the definition for the resource", configv1.Secret{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/secrets/{name}").To(u.findTeamSecret).
			Operation("GetTeamSecret").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the secert in the name")).
			Doc("Used to retrieve the secret from the team").
			Returns(http.StatusOK, "Contains the definition for the resource", configv1.Secret{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/secrets/{name}").To(u.updateTeamSecret).
			Operation("UpdateTeamSecret").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of secret you are creating / updating")).
			Doc("Used to update the secret in the team").
			Reads(configv1.Secret{}, "The definition for the secret you are creating or updating").
			Returns(http.StatusOK, "Contains updated definition of the secret", configv1.Secret{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/secrets/{name}").To(u.deleteTeamSecret).
			Operation("DeleteTeamSecret").
			Param(ws.PathParameter("name", "Is name the of the secret you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to delete the secret from team").
			Returns(http.StatusOK, "Contains the former definition of the secret", configv1.Secret{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)
}

// updateTeamSecret is used to add a secret to the team
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
