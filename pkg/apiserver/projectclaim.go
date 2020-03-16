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

/*
// findProjectClaims returns a list of credential claims
func (u teamHandler) findProjectClaims(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).Cloud().GCP().Projects().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findProjectClaims returns a specific credential
func (u teamHandler) findProjectClaim(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).Cloud().GCP().Projects().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateProjectClaim is used to update an credential claim for a team
func (u teamHandler) updateProjectClaim(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &gcp.ProjectClaims{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).Cloud().GCP().Projects().Update(req.Request.Context(), claim)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteProjectClaim is used to remove a credential from a team cluster
func (u teamHandler) deleteProjectClaim(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).ProjectClaims().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}
*/
