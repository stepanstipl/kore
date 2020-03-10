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
	"net/http"

	"github.com/appvia/kore/pkg/kore/authentication"

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
