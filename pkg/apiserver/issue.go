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

package apiserver

import (
	"net/http"
	"time"

	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&issueHandler{})
}

type issueHandler struct {
	kore.Interface
	// default handler
	DefaultHandler
}

// Path returns the handler path
func (l *issueHandler) Path() string {
	return "issue"
}

// Register is responsible for handling the registration
func (l *issueHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add(l.Path())
	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the issue webservice with container")

	l.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		ws.PUT("/authorize").To(l.issueHandler).
			Filter(filters.NewRateLimiter(filters.RateConfig{Period: 60 * time.Second, Limit: 5})).
			Doc("Used to auhorize and swap tokens for a locally minted token").
			Operation("LocalAuthorize").
			Param(ws.QueryParameter("method", "The type of token being swapped, i.e basicauth, openid").Required(true)).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// issueHandler is responsible for issuing a local token for local users
func (l *issueHandler) issueHandler(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		// @step: get the token type being swapped
		method := req.QueryParameter("method")
		switch method {
		case "basicauth":
			issued, err := l.Users().Identities().IssueToken(req.Request.Context(), "kubernetes", []string{"impersonate"})
			if err != nil {
				return err
			}
			return resp.WriteHeaderAndEntity(http.StatusOK, &types.IssuedToken{Token: issued})

		default:
			resp.WriteHeader(http.StatusNotImplemented)
		}

		return nil
	})
}

// EnableAdminsOnly indicates if we need to be an admin user
func (l *issueHandler) EnableAdminsOnly() bool {
	return false
}

// Name returns the name of the handler
func (l *issueHandler) Name() string {
	return "issue"
}
