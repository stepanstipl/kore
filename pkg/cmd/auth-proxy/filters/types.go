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

	"github.com/julienschmidt/httprouter"
)

// Middleware is used to define the middleware
type Middleware interface {
	Serve(http.Handler) http.Handler
}

// Interface is the contract to the middleware filters
type Interface interface {
	// Wrap is used to call the chain and handler
	Wrap(httprouter.Handle) httprouter.Handle
	// Use appends a middleware to the chain
	Use(...Middleware)
}
