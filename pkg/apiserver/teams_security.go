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

	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addTeamSecurityRoutes(ws *restful.WebService) {
	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/security")).To(u.getTeamSecurityOverview).
			Doc("Returns an overview of the security posture for resources owned by this team").
			Operation("GetTeamSecurityOverview").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "The requested security overview", securityv1.SecurityOverview{}),
	)
}

func (t *teamHandler) getTeamSecurityOverview(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		overview, err := t.Security().GetTeamOverview(req.Request.Context(), req.PathParameter("team"))
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, overview)
	})
}
