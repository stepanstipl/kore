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
	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
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
		withAllErrors(ws.POST("cluster")).To(c.estimateClusterPlanCost).
			Doc("Returns the estimated cost of the supplied cluster plan").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("EstimateClusterPlanCost").
			Reads(configv1.Plan{}, "The specification for the plan you want estimating").
			Returns(http.StatusOK, "An estimate of the costs for the cluster plan", costsv1.CostEstimate{}),
	)

	ws.Route(
		withAllErrors(ws.POST("service")).To(c.estimateServicePlanCost).
			Doc("Returns the estimated cost of the supplied service plan").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("EstimateServicePlanCost").
			Reads(servicesv1.ServicePlan{}, "The specification for the plan you want estimating").
			Returns(http.StatusOK, "An estimate of the costs for the service plan", costsv1.CostEstimate{}),
	)

	return ws, nil
}

func (c costEstimatesHandler) estimateClusterPlanCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan := &configv1.Plan{}
		if err := req.ReadEntity(plan); err != nil {
			return err
		}

		estimate, err := c.Costs().Estimates().GetClusterEstimate(&plan.Spec)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, estimate)
	})
}

func (c costEstimatesHandler) estimateServicePlanCost(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan := &servicesv1.ServicePlan{}
		if err := req.ReadEntity(plan); err != nil {
			return err
		}

		estimate, err := c.Costs().Estimates().GetServiceEstimate(&plan.Spec)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, estimate)
	})
}

// Name returns the name of the handler
func (c costEstimatesHandler) Name() string {
	return "costestimates"
}
