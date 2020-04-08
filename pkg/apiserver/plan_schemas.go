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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&planSchemasHandler{})
}

type planSchemasHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *planSchemasHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("planschemas")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the planschemas webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.getPlanSchema).
			Doc("Returns a specific plan schema from the kore").
			Operation("GetPlanSchema").
			Param(ws.PathParameter("name", "The name of the plan schema you wish to retrieve")).
			Returns(http.StatusOK, "Contains the plan schema definition from the kore", configv1.PlanPolicy{}),
	)

	return ws, nil
}

func (p planSchemasHandler) getPlanSchema(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		schema := ""
		switch req.PathParameter("name") {
		case "GKE":
			schema = assets.GKEPlanSchema
			break
		case "EKS":
			schema = assets.EKSPlanSchema
			break
		default:
			return resp.WriteHeaderAndEntity(http.StatusNotFound, nil)
		}
		// This is ALWAYS json as we're returning json schema, so don't want
		// to use normal content type negotiation.
		resp.AddHeader("Content-Type", restful.MIME_JSON)
		resp.WriteHeader(http.StatusOK)
		if _, err := resp.Write([]byte(schema)); err != nil {
			return err
		}
		return nil
	})
}

// Name returns the name of the handler
func (p planSchemasHandler) Name() string {
	return "planschemas"
}
