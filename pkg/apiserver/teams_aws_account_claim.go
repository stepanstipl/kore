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

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) addAWSAccountClaimRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/awsaccountclaims").To(u.findAWSAccountClaims).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used to return a list of AWS accounts which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSAccountClaimList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/awsaccountclaims/{name}").To(u.findAWSAccountClaim).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a list of AWS Accounts which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSAccountClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findAWSAccountClaims returns a list of credential claims
func (u teamHandler) findAWSAccountClaims(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().AWS().AWSAccountClaims().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findAWSAccountClaim returns a specific credential
func (u teamHandler) findAWSAccountClaim(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).Cloud().AWS().AWSAccountClaims().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}
