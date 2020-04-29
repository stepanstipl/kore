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
	"strings"
	"time"

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"

	restful "github.com/emicklei/go-restful"
)

// NewAuditingFilter creates a new Auditing for a route and returns the filter function.
func NewAuditingFilter(audit func() kore.Audit, apiVersion string, resource string, operation string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		user := req.Request.Context().Value(authentication.ContextKey{})

		var username string
		if user != nil {
			username = user.(authentication.Identity).Username()
		}

		// Normalise resource URI to strip the version - it may or may not have
		// a trailing slash so trim a possible trailing slash off for consistency.
		uri := strings.TrimSuffix(strings.TrimPrefix(req.Request.URL.EscapedPath(), apiVersion), "/")
		// Remove the API version prefix from the resource if present.
		resource := strings.TrimPrefix(resource, apiVersion)

		// Default to internal error - any non-internal-error case will change this
		// after running the chain.
		responseCode := 500
		start := time.Now()

		defer func() {
			finish := time.Now()

			// @TODO: Refine what is audited with a policy in future, for now, just audit everything.
			audit().Record(req.Request.Context(),
				persistence.Resource(resource),
				persistence.ResourceURI(uri),
				persistence.APIVersion(apiVersion),
				persistence.Verb(req.Request.Method),
				persistence.Operation(operation),
				persistence.Team(req.PathParameter("team")), // Might be nil for some paths, but that's OK.
				persistence.User(username),
				persistence.StartedAt(start),
				persistence.CompletedAt(finish),
				persistence.ResponseCode(responseCode),
			).Event(fmt.Sprintf("%s: %s %s", operation, req.Request.Method, req.Request.URL.EscapedPath()))
		}()

		chain.ProcessFilter(req, resp)
		responseCode = resp.StatusCode()
	}
}
