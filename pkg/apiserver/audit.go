/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
		ws.GET("").To(a.findTeamsAudit).
			Doc("Used to return all the audit event across all the teams").
			Param(ws.QueryParameter("since", "The time duration to return the events within").DefaultValue("60m")).
			Returns(http.StatusOK, "A collection of events from the team", orgv1.AuditEventList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
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
