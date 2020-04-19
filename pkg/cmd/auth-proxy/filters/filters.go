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
	"sync"

	"github.com/julienschmidt/httprouter"
)

type filtersImpl struct {
	sync.RWMutex
	// handlers is a collection of middleware
	handlers []Middleware
}

// New creates and returns a filters wrapper
func New() Interface {
	return &filtersImpl{}
}

// Wrap is used to call the chain and handler
func (f *filtersImpl) Wrap(handler httprouter.Handle) httprouter.Handle {
	size := len(f.handlers)
	if size == 0 {
		return handler
	}

	// we use the above handler as the last thing to call
	v := f.handlers[size-1].Serve(wrapHandler(handler))

	for i := 0; i < (size - 1); i++ {
		v = f.handlers[size-(2+i)].Serve(v)
	}

	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		v.ServeHTTP(w, req)
	}
}

// Use appends a middleware to the chain
func (f *filtersImpl) Use(handlers ...Middleware) {
	f.Lock()
	defer f.Unlock()

	f.handlers = append(f.handlers, handlers...)
}

func wrapHandler(handle httprouter.Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if handle != nil {
			handle(w, req, nil)
		}
	})
}
