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
	"fmt"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/kore/authentication"

	"github.com/appvia/kore/pkg/apiserver/filters"

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
	// DefaultHandler implements default features
	DefaultHandler
}

func (p *servicePlansHandler) systemServicePlanFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		servicePlan, err := p.ServicePlans().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if servicePlan != nil && servicePlan.Annotations[kore.AnnotationSystem] == "true" {
			resp.WriteHeader(http.StatusForbidden)
			return nil
		}

		// @step: continue with the chain
		chain.ProcessFilter(req, resp)
		return nil
	})
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
		withAllNonValidationErrors(ws.GET("")).To(p.listServicePlans).
			Doc("Returns all the available service plans").
			Operation("ListServicePlans").
			Param(ws.QueryParameter("kind", "Filters service plans for a specific kind")).
			Returns(http.StatusOK, "A list of service plans", servicesv1.ServicePlanList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.getServicePlan).
			Doc("Returns a specific service plan").
			Operation("GetServicePlan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to retrieve")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}/schema")).To(p.getServicePlanSchema).
			Doc("Returns the JSON schema for the plan. If a plan doesn't have a schema, it returns the JSON schema defined on the service kind").
			Operation("GetServicePlanSchema").
			Param(ws.PathParameter("name", "The name of the service plan")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service schema definition", map[string]interface{}{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}/credentialschema")).To(p.getServiceCredentialSchema).
			Doc("Returns the JSON schema for the service credentials defined in the plan. If a plan doesn't have credential schema, it returns the JSON schema defined on the service kind").
			Operation("GetServiceCredentialSchema").
			Param(ws.PathParameter("name", "The name of the service plan")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service credential schema definition", map[string]interface{}{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updateServicePlan).
			Filter(filters.Admin).
			Filter(p.systemServicePlanFilter).
			Doc("Creates or updates a service plan").
			Operation("UpdateServicePlan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to create or update")).
			Reads(servicesv1.ServicePlan{}, "The specification for the service plan you are creating or updating").
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deleteServicePlan).
			Filter(filters.Admin).
			Filter(p.systemServicePlanFilter).
			Doc("Deletes a service plan").
			Operation("DeleteServicePLan").
			Param(ws.PathParameter("name", "The name of the service plan you wish to delete")).
			Returns(http.StatusNotFound, "the service plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service plan definition", servicesv1.ServicePlan{}),
	)

	return ws, nil
}

// getServicePlan returns a specific service plan
func (p servicePlansHandler) getServicePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan, err := p.ServicePlans().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// getServicePlan returns the schema for the given service plan
func (p servicePlansHandler) getServicePlanSchema(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		schema, err := p.ServicePlans().GetSchema(req.Request.Context(), req.PathParameter("name"))
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

// getServiceCredentialSchema returns the credential schema for the given service plan
func (p servicePlansHandler) getServiceCredentialSchema(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		schema, err := p.ServicePlans().GetCredentialSchema(req.Request.Context(), req.PathParameter("name"))
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

// listServicePlans returns all service plans in the kore
func (p servicePlansHandler) listServicePlans(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		user := authentication.MustGetIdentity(req.Request.Context())
		kind := strings.ToLower(req.QueryParameter("kind"))

		list, err := p.ServicePlans().ListFiltered(req.Request.Context(), func(plan servicesv1.ServicePlan) bool {
			if kind != "" && plan.Kind != kind {
				return false
			}
			if !user.IsGlobalAdmin() && plan.Annotations[kore.AnnotationSystem] == "true" {
				return false
			}

			return true
		})
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
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

		if plan.Annotations[kore.AnnotationSystem] != "" {
			writeError(req, resp, fmt.Errorf("setting %q annotation is not allowed", kore.AnnotationSystem), http.StatusForbidden)
			return nil
		}

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
