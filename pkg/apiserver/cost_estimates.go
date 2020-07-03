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
	RegisterHandler(&costEstimatesHandler{})
}

type costEstimatesHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (c *costEstimatesHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("costestimates")
	tags := []string{"costestimates"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the costs estimation webservice")

	c.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("/metadata/{cloud}/regions")).To(c.getMetadataRegions).
			Doc("Returns regions").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetMetadataRegions").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve regions for")).
			Returns(http.StatusNotFound, "cloud doesn't exist", nil).
			Returns(http.StatusOK, "A list of all the regions organised by continent", costsv1.ContinentList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/metadata/{cloud}/regions/{region}/instances")).To(c.getMetadataInstances).
			Doc("Returns prices and instance types for a given region of a given cloud provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetMetadataInstances").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve instance types/prices for")).
			Param(ws.PathParameter("region", "The region to retrieve instance types/prices for")).
			Returns(http.StatusNotFound, "cloud or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of instance types with their pricing", costsv1.InstanceTypeList{}),
	)

	ws.Route(
		withAllErrors(ws.POST("cluster")).To(c.estimateClusterPlanCost).
			Doc("Returns the estimated cost of the supplied cluster plan").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("EstimateClusterPlanCost").
			Returns(http.StatusOK, "An estimate of the costs for the cluster plan", costsv1.CostEstimate{}),
	)

	ws.Route(
		withAllErrors(ws.POST("service")).To(c.estimateServicePlanCost).
			Doc("Returns the estimated cost of the supplied service plan").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("EstimateServicePlanCost").
			Returns(http.StatusOK, "An estimate of the costs for the service plan", costsv1.CostEstimate{}),
	)

	return ws, nil
}

func (c costEstimatesHandler) getMetadataRegions(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.ContinentList{})
	})
}

func (c costEstimatesHandler) getMetadataInstances(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.InstanceTypeList{})
	})
}

func (c costEstimatesHandler) estimateClusterPlanCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.CostEstimate{})
	})
}

func (c costEstimatesHandler) estimateServicePlanCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		return errors.New("not implemented")
		// return resp.WriteHeaderAndEntity(http.StatusOK, &costsv1.CostEstimate{})
	})
}

// Name returns the name of the handler
func (c costEstimatesHandler) Name() string {
	return "costestimates"
}
