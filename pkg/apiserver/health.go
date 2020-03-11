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
	log.Info("registering the health webservice with container")

	l.Interface = i

	ws := &restful.WebService{}

	ws.Route(
		ws.GET("/healthz").To(l.healthHandler).
			Doc("Used to start the authorization flow for user authentication").
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
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
