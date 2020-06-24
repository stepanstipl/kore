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

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&costsHandler{})
}

type costsHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (c *costsHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("costs")
	tags := []string{"security"}

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
			Doc("Returns all the costs").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListCosts").
			Returns(http.StatusOK, "A list of all the costs", costsv1.CostList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/metadata/{cloud}/regions")).To(c.getRegions).
			Doc("Returns regions").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetRegions").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve regions for")).
			Returns(http.StatusNotFound, "cloud doesn't exist", nil).
			Returns(http.StatusOK, "A list of all the regions organised by ", costsv1.InstanceTypeList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/metadata/{kind}/regions/{region}/instances")).To(c.getPrices).
			Doc("Returns prices and instance types for a given region of a given cloud provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetPrices").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve instance types for")).
			Param(ws.PathParameter("region", "The region to retrieve instance types for")).
			Returns(http.StatusNotFound, "cloud or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of all the costs", costsv1.InstanceTypeList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(c.getCost).
			Doc("Returns a specific cost from the kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCost").
			Param(ws.PathParameter("name", "The name of the cost you wish to retrieve")).
			Returns(http.StatusNotFound, "the cost with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the cost definition", costsv1.Cost{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(c.updateCost).
			Doc("Used to create or update a cost in the kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("UpdateCost").
			Param(ws.PathParameter("name", "The name of the cost you wish to create or update")).
			Reads(costsv1.Cost{}, "The specification for the cost you are creating or updating").
			Returns(http.StatusOK, "Contains the cost definition", costsv1.Cost{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(c.deleteCost).
			Doc("Used to delete a cost from the kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("RemoveCost").
			Param(ws.PathParameter("name", "The name of the cost you wish to delete")).
			Returns(http.StatusNotFound, "the cost with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the cost definition", costsv1.Cost{}),
	)

	return ws, nil
}

// getCost returns a specific cost
func (c costsHandler) getCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan, err := c.Costs().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// getRegions returns the regions metadata
func (c costsHandler) getRegions(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.ContinentList{})
	})
}

// getPrices returns the pricing metadata
func (c costsHandler) getPrices(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.InstanceTypeList{})
	})
}

// listCosts returns all costs in the kore
func (c costsHandler) listCosts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		costs, err := c.Costs().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, costs)
	})
}

// updateCost is used to update or create a cost in the kore
func (c costsHandler) updateCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		cost := &costsv1.Cost{}
		if err := req.ReadEntity(cost); err != nil {
			return err
		}
		cost.Name = name

		if err := c.Costs().Update(req.Request.Context(), cost, false); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, cost)
	})
}

// deleteCost is used to update or create a cost in the kore
func (c costsHandler) deleteCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		cost, err := c.Costs().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, cost)
	})
}

// Name returns the name of the handler
func (c costsHandler) Name() string {
	return "costs"
}
