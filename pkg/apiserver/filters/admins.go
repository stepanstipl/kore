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
package filters

import (
	"net/http"

	"github.com/appvia/kore/pkg/hub/authentication"

	restful "github.com/emicklei/go-restful"
)

var (
	// DefaultAdminHandler is the default filter
	DefaultAdminHandler AdminHandler
)

// AdminHandler is the default authentication handler
type AdminHandler struct{}

// Filter is called on middleware execution
func (a *AdminHandler) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// @step: ensure the user is a member of the team
	user := authentication.MustGetIdentity(req.Request.Context())
	if !user.IsGlobalAdmin() {
		resp.WriteHeader(http.StatusForbidden)

		return
	}

	// @step: continue with the chain
	chain.ProcessFilter(req, resp)
}
