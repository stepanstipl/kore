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
	"time"

	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/services/users"

	restful "github.com/emicklei/go-restful"
)

// Audit represents a function to retrieve an implementation of the Audit API.
type Audit func() users.Audit

// Auditing provides auditing facilities for callers of the API.
type Auditing struct {
	Audit     Audit
	Resource  string
	Operation string
}

// NewAuditingFilter creates a new Auditing for a route and returns the filter function.
func NewAuditingFilter(audit Audit, resource string, operation string) restful.FilterFunction {
	a := Auditing{
		Audit:     audit,
		Resource:  resource,
		Operation: operation,
	}
	return a.Filter
}

// Filter implements the auditing filter for the api server
func (a *Auditing) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	user := req.Request.Context().Value(authentication.ContextKey{})

	var username string
	if user != nil {
		username = user.(authentication.Identity).Username()
	}

	// Default to internal error - any non-internal-error case will change this
	// after running the chain.
	result := 500
	start := time.Now()

	defer func() {
		finish := time.Now()

		a.Audit().Record(req.Request.Context(),
			users.Resource(a.Resource),
			users.ResourceURI(req.Request.URL.EscapedPath()),
			users.Verb(req.Request.Method),
			users.Operation(a.Operation),
			users.Team(req.PathParameter("team")), // Might be nil for some paths, but that's OK.
			users.User(username),
			users.StartedAt(start),
			users.CompletedAt(finish),
			users.Result(result),
		).Event(fmt.Sprintf("%s: %s %s", a.Operation, req.Request.Method, req.Request.URL.EscapedPath()))
	}()

	chain.ProcessFilter(req, resp)
	result = resp.StatusCode()
}
