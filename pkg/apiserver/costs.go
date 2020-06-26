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
	"errors"
	"net/http"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&costHandler{})
}

type costHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (c *costHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("costs")
	tags := []string{"costs"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the costs webservice")

	c.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(c.listCosts).
			Doc("Returns a list of actual costs").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListCosts").
			Returns(http.StatusOK, "A list of all costs known to the system", costsv1.CostList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(c.getCost).
			Doc("Gets a specific cost").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCost").
			Param(ws.PathParameter("name", "The name of the cost to retrieve")).
			Returns(http.StatusNotFound, "Cost doesn't exist", nil).
			Returns(http.StatusOK, "Cost found", costsv1.Cost{}),
	)

	return ws, nil
}

func (c costHandler) listCosts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.CostList{})
	})
}

func (c costHandler) getCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.Cost{})
	})
}

// Name returns the name of the handler
func (c costHandler) Name() string {
	return "costs"
}
