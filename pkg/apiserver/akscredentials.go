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

	aks "github.com/appvia/kore/pkg/apis/aks/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

// listAKSCredentials returns all the AKS credentials under the team
func (u teamHandler) listAKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().AKSCredentials().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getAKSCredentials returns the named AKS credentials under the team
func (u teamHandler) getAKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().AKSCredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// deleteAKSCredentials is responsible for deleting the AKS credentials
func (u teamHandler) deleteAKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Cloud().AKSCredentials().Get(ctx, name)
		if err != nil {
			return err
		}

		err = u.Teams().Team(team).Cloud().AKSCredentials().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateAKSCredentials is responsible for creating/updating AKS credentials
func (u teamHandler) updateAKSCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &aks.AKSCredentials{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if _, err := u.Teams().Team(team).Cloud().AKSCredentials().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
