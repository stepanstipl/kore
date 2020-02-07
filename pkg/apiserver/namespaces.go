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

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"

	restful "github.com/emicklei/go-restful"
)

// findNamespaces returns a list of namespace claims
func (u teamHandler) findNamespaces(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).NamespaceClaims().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findNamespace returns a specific namespace
func (u teamHandler) findNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).NamespaceClaims().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateNamespace is used to update an namespace claim for a team
func (u teamHandler) updateNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &clustersv1.NamespaceClaim{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).NamespaceClaims().Update(req.Request.Context(), claim)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteNamespace is used to remove a namespace from a team cluster
func (u teamHandler) deleteNamespace(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).NamespaceClaims().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}
