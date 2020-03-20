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
