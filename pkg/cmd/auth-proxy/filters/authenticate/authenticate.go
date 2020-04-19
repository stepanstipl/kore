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

package authenticate

import (
	"errors"
	"net/http"

	"github.com/appvia/kore/pkg/cmd/auth-proxy/filters"
	"github.com/appvia/kore/pkg/cmd/auth-proxy/verifiers"
)

// Options are configurable for the filter
type Options struct {
	Verifiers []verifiers.Interface
}

type authImpl struct {
	Options
}

// New creates and return the filter
func New(options Options) (filters.Middleware, error) {
	if len(options.Verifiers) == 0 {
		return nil, errors.New("no verifiers")
	}

	return &authImpl{Options: options}, nil
}

// Serve handles the filter
func (a *authImpl) Serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		allowed := func() bool {
			for _, x := range a.Verifiers {
				matched, err := x.Admit(req)
				if err != nil {
					continue
				}
				if matched {
					return true
				}
			}

			return false
		}()
		if !allowed {
			authFailureCounter.Inc()

			w.WriteHeader(http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, req)
	})
}
