/*
 * Copyright (C) 2019 Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package filters

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/appvia/kore/pkg/apiserver/plugins/identity"
	"github.com/appvia/kore/pkg/hub/authentication"

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
