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

	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&servicePlanSchemasHandler{})
}

type servicePlanSchemasHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *servicePlanSchemasHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("serviceplanschemas")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the serviceplanschemas webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.getServicePlanSchema).
			Doc("Returns a specific service plan schema from the kore").
			Operation("GetServicePlanSchema").
			Param(ws.PathParameter("name", "The name of the service plan schema you wish to retrieve")).
			Returns(http.StatusOK, "Contains the service plan schema definition", nil),
	)

	return ws, nil
}

func (p servicePlanSchemasHandler) getServicePlanSchema(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		kind := req.PathParameter("name")
		provider := p.ServiceProviders().GetProviderForKind(kind)
		if provider == nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, nil)
		}

		schema, err := provider.JSONSchema(kind, "")
		if err != nil {
			return err
		}

		resp.AddHeader("Content-Type", restful.MIME_JSON)
		resp.WriteHeader(http.StatusOK)
		if _, err := resp.Write([]byte(schema)); err != nil {
			return err
		}
		return nil
	})
}

// Name returns the name of the handler
func (p servicePlanSchemasHandler) Name() string {
	return "serviceplanschemas"
}
