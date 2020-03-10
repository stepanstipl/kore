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
	"context"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/kore/authentication"

	restful "github.com/emicklei/go-restful"
)

var (
	// DefaultAuthentication is the default filter
	DefaultAuthentication AuthenticationHandler
)

// AuthenticationHandler is the default authentication handler
type AuthenticationHandler struct {
	// Realm specify a authentication realm to if any
	Realm string
}

// Filter is called on middleware execution
func (a *AuthenticationHandler) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if len(identity.GetPlugins()) <= 0 {
		chain.ProcessFilter(req, resp)

		return
	}

	// @step: find a matching ident
	user, found := func() (authentication.Identity, bool) {
		rq := &requestorImpl{request: req}

		for _, p := range identity.GetPlugins() {
			if u, found := p.Admit(req.Request.Context(), rq); found {
				return u, found
			}
		}

		return nil, false
	}()
	if !found {

		if a.Realm != "" {
			resp.Header().Set("WWW-Authenticate", fmt.Sprintf("realm=\"%s\"", a.Realm))
		}
		resp.WriteHeader(http.StatusUnauthorized)

		return
	}
	ctx := context.WithValue(req.Request.Context(), authentication.ContextKey{}, user)

	// @step: add the user into the context
	req.Request = req.Request.WithContext(ctx)

	// @step: continue with the chain
	chain.ProcessFilter(req, resp)
}

type requestorImpl struct {
	request *restful.Request
}

func (r requestorImpl) ClientCertficate() *x509.Certificate {
	if r.request.Request.TLS != nil {
		return r.request.Request.TLS.PeerCertificates[0]
	}

	return nil
}

func (r requestorImpl) Headers() http.Header {
	return r.request.Request.Header
}
