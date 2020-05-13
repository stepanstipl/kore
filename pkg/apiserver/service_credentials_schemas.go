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
	RegisterHandler(&serviceCredentialSchemasHandler{})
}

type serviceCredentialSchemasHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *serviceCredentialSchemasHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("servicecredentialschemas")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the servicecredentialschemas webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{kind}")).To(p.getServiceCredentialSchema).
			Doc("Returns a specific service credential schema").
			Operation("GetServiceCredentialSchemaForKind").
			Param(ws.PathParameter("kind", "The service kind")).
			Returns(http.StatusOK, "Contains the service credential schema definition", nil),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{kind}/{name}")).To(p.getServiceCredentialSchema).
			Doc("Returns a specific service credential schema").
			Operation("GetServiceCredentialSchemaForPlan").
			Param(ws.PathParameter("kind", "The service kind")).
			Param(ws.PathParameter("name", "The service plan name")).
			Returns(http.StatusOK, "Contains the service credential schema definition", nil),
	)

	return ws, nil
}

func (p serviceCredentialSchemasHandler) getServiceCredentialSchema(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		var planName string

		kind := req.PathParameter("kind")
		name := req.PathParameter("name")
		if name != "" {
			plan, err := p.ServicePlans().Get(req.Request.Context(), name)
			if err != nil {
				return err
			}
			planName = plan.PlanShortName()
		}

		provider := p.ServiceProviders().GetProviderForKind(kind)
		if provider == nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, nil)
		}

		schema, err := provider.CredentialsJSONSchema(kind, planName)
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
func (p serviceCredentialSchemasHandler) Name() string {
	return "servicecredentialschemas"
}
