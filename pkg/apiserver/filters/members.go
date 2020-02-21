/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package filters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
)

var (
	// DefaultMembersHandler is the default filter
	DefaultMembersHandler MembersHandler
)

// MembersHandler is the default authentication handler
type MembersHandler struct {
}

// Filter is called on middleware execution
func (a *MembersHandler) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	team := req.PathParameter("team")
	if team == "" {
		// @step: continue with the chain
		chain.ProcessFilter(req, resp)

		return
	}

	// @step: ensure the user is a member of the team
	user := authentication.MustGetIdentity(req.Request.Context())
	if !user.IsGlobalAdmin() {
		// @TODO this needs to be removed - its VERY HACKY - but we will keep until
		// we have a proper review of the API
		if !strings.HasSuffix(req.Request.RequestURI, fmt.Sprintf("teams/%s", team)) {
			if !utils.Contains(team, user.Teams()) {
				resp.WriteHeader(http.StatusForbidden)

				return
			}
		}
	}

	// @step: continue with the chain
	chain.ProcessFilter(req, resp)
}
