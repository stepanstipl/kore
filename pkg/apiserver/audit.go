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

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&auditHandler{})
}

type auditHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (a *auditHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("audit")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the audit webservice")

	a.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllErrors(ws.GET("")).To(a.findTeamsAudit).
			Doc("Used to return all the audit event across all the teams").
			Operation("ListAuditEvents").
			Param(ws.QueryParameter("since", "The time duration to return the events within").DefaultValue("60m")).
			Returns(http.StatusOK, "A collection of events from the team", orgv1.AuditEventList{}),
	)

	return ws, nil
}

// findTeamsAudit returns all the audit events across all the teams
func (a *auditHandler) findTeamsAudit(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		since := req.QueryParameter("since")
		if since == "" {
			since = "60m"
		}
		tm, err := time.ParseDuration(since)
		if err != nil {
			return err
		}

		list, err := a.Teams().AuditEvents(req.Request.Context(), tm)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// Name returns the name of the handler
func (a auditHandler) Name() string {
	return "audit"
}
