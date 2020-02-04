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
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&healthHandler{})
}

type healthHandler struct {
	kore.Interface
	// provides the default handler
	DefaultHandler
}

// Register is responsible for handling the registration
func (l *healthHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.Info("registering the healt webservice with container")

	l.Interface = i

	ws := &restful.WebService{}

	ws.Route(
		ws.GET("/healthz").To(l.healthHandler).
			Doc("Used to start the authorization flow for user authentication").
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// healthHandler is responsible for authorizing a client
func (l *healthHandler) healthHandler(req *restful.Request, resp *restful.Response) {
	_, _ = resp.Write([]byte("OK"))
}

// EnableAuthentication indicates if this service needs auth
func (l *healthHandler) EnableAuthentication() bool {
	return false
}

// EnableLogging indicates if logging is one
func (l *healthHandler) EnableLogging() bool {
	return false
}

// Name returns the name of the handler
func (l *healthHandler) Name() string {
	return "health"
}
