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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	RegisterHandler(&plansHandler{})
}

type plansHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *plansHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("plans")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the plans webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(p.findPlans).
			Doc("Returns all the classes available to initialized in the kore").
			Operation("ListPlans").
			Param(ws.QueryParameter("kind", "Returns all plans for a specific resource type")).
			Returns(http.StatusOK, "A list of all the plans", configv1.PlanList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.findPlan).
			Doc("Returns a specific class plan from the kore").
			Operation("GetPlan").
			Param(ws.PathParameter("name", "The name of the plan you wish to retrieve")).
			Returns(http.StatusNotFound, "the plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the plan definition", configv1.Plan{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updatePlan).
			Doc("Used to create or update a plan in the kore").
			Operation("UpdatePlan").
			Param(ws.PathParameter("name", "The name of the plan you wish to create or update")).
			Reads(configv1.Plan{}, "The specification for the plan you are creating or updating").
			Returns(http.StatusOK, "Contains the plan definition", configv1.Plan{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deletePlan).
			Doc("Used to delete a plan from the kore").
			Operation("RemovePlan").
			Param(ws.PathParameter("name", "The name of the plan you wish to delete")).
			Returns(http.StatusNotFound, "the plan with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the plan definition", configv1.Plan{}),
	)

	return ws, nil
}

// findPlan returns a specific plan
func (p plansHandler) findPlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan, err := p.Plans().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// findPlans returns all plans in the kore
func (p plansHandler) findPlans(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plans, err := p.Plans().List(req.Request.Context())
		if err != nil {
			return err
		}

		kind := strings.ToLower(req.QueryParameter("kind"))

		filtered := &configv1.PlanList{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PlanList",
			},
			Items: []configv1.Plan{},
		}
		for _, x := range plans.Items {
			if kind != "" && strings.ToLower(x.Spec.Kind) != kind {
				continue
			}
			filtered.Items = append(filtered.Items, x)
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, filtered)
	})
}

// updatePlan is used to update or create a plan in the kore
func (p plansHandler) updatePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		plan := &configv1.Plan{}
		if err := req.ReadEntity(plan); err != nil {
			return err
		}
		plan.Name = name

		if err := p.Plans().Update(req.Request.Context(), plan, false); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// deletePlan is used to update or create a plan in the kore
func (p plansHandler) deletePlan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		plan, err := p.Plans().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, plan)
	})
}

// Name returns the name of the handler
func (p plansHandler) Name() string {
	return "plans"
}
