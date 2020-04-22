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
	"strings"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&servicePlansHandler{})
}

type servicePlansHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *servicePlansHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("serviceplans")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the serviceplans webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(p.findServicePlans).
			Doc("Returns all the available service plans").
			Operation("ListServicePlans").
			Param(ws.QueryParameter("kind", "Filters service plans for a specific kind")).
			Returns(http.StatusOK, "A list of service plans", servicesv1.ServicePlanList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.findServicePlan).
			Doc("Returns a specific service plan").
			Operation("GetServicePlan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to retrieve")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updateServicePlan).
			Doc("Creates or updates a service plan").
			Operation("UpdateServicePlan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to create or update")).
			Reads(servicesv1.ServicePlan{}, "The specification for the service plan you are creating or updating").
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deleteServicePlan).
			Doc("Deletes a service plan").
			Operation("DeleteServicePLan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to delete")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	return ws, nil
}

// findServicePlan returns a specific service plan
func (p servicePlansHandler) findServicePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan, err := p.ServicePlans().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// findServicePlans returns all service plans in the kore
func (p servicePlansHandler) findServicePlans(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		res, err := p.ServicePlans().List(req.Request.Context())
		if err != nil {
			return err
		}

		kind := strings.ToLower(req.QueryParameter("kind"))
		if kind != "" {
			var items []servicesv1.ServicePlan
			for _, x := range res.Items {
				if strings.ToLower(x.Spec.Kind) == kind {
					items = append(items, x)
				}
			}
			res.Items = items
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, res)
	})
}

// updateServicePlan is used to update or create a service plan in the kore
func (p servicePlansHandler) updateServicePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		plan := &servicesv1.ServicePlan{}
		if err := req.ReadEntity(plan); err != nil {
			return err
		}
		plan.Name = name

		if err := p.ServicePlans().Update(req.Request.Context(), plan); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// deleteServicePlan is used to update or create a service plan in the kore
func (p servicePlansHandler) deleteServicePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		plan, err := p.ServicePlans().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// Name returns the name of the handler
func (p servicePlansHandler) Name() string {
	return "serviceplans"
}

// Enabled returns true if the services feature gate is enabled
func (p servicePlansHandler) Enabled() bool {
	return p.Config().IsFeatureGateEnabled(kore.FeatureGateServices)
}
