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

	restful "github.com/emicklei/go-restful"
)

var (
	// DefaultNotImplementedHandler is the default filter
	DefaultNotImplementedHandler NotImplementedHandler
)

// NotImplementedHandler is the default authentication handler
type NotImplementedHandler struct{}

// Filter is called on middleware execution
func (a *NotImplementedHandler) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.WriteHeader(http.StatusNotImplemented)
}
