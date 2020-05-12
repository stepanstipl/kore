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
	RegisterHandler(&planPoliciesHandler{})
}

type planPoliciesHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *planPoliciesHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("planpolicies")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the planpolicies webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		ws.GET("").To(p.findPlanPolicies).
			Doc("Returns all the plan policies").
			Operation("ListPlanPolicies").
			Param(ws.QueryParameter("kind", "Returns all plan policies for a specific resource type")).
			Returns(http.StatusOK, "A list of all the plan policies in the kore", configv1.PlanPolicyList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{name}").To(p.findPlanPolicy).
			Doc("Returns a specific plan policy from the kore").
			Operation("GetPlanPolicy").
			Param(ws.PathParameter("name", "The name of the plan policy you wish to retrieve")).
			Returns(http.StatusOK, "Contains the plan policy definition from the kore", configv1.PlanPolicy{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{name}").To(p.updatePlanPolicy).
			Doc("Used to create or update a plan policy in the kore").
			Operation("UpdatePlanPolicy").
			Param(ws.PathParameter("name", "The name of the plan policy you wish to update")).
			Reads(configv1.PlanPolicy{}, "The specification for the plan policy you are updating").
			Returns(http.StatusOK, "Contains the plan policy definition from the kore", configv1.PlanPolicy{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{name}").To(p.deletePlanPolicy).
			Doc("Used to delete a plan policy from the kore").
			Operation("RemovePlanPolicy").
			Param(ws.PathParameter("name", "The name of the plan policy you wish to delete")).
			Returns(http.StatusOK, "Contains the plan policy definition from the kore", configv1.PlanPolicy{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

func (p planPoliciesHandler) findPlanPolicy(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		planPolicy, err := p.PlanPolicies().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, planPolicy)
	})
}

func (p planPoliciesHandler) findPlanPolicies(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		planPolicies, err := p.PlanPolicies().List(req.Request.Context())
		if err != nil {
			return err
		}

		kind := strings.ToLower(req.QueryParameter("kind"))

		filtered := &configv1.PlanPolicyList{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PlanPolicyList",
			},
			Items: []configv1.PlanPolicy{},
		}
		for _, x := range planPolicies.Items {
			if kind != "" && strings.ToLower(x.Spec.Kind) != kind {
				continue
			}
			filtered.Items = append(filtered.Items, x)
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, filtered)
	})
}

func (p planPoliciesHandler) updatePlanPolicy(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		planPolicy := &configv1.PlanPolicy{}
		if err := req.ReadEntity(planPolicy); err != nil {
			return err
		}
		planPolicy.Name = name

		if err := p.PlanPolicies().Update(req.Request.Context(), planPolicy, false); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, planPolicy)
	})
}

func (p planPoliciesHandler) deletePlanPolicy(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		planPolicy, err := p.PlanPolicies().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, planPolicy)
	})
}

// Name returns the name of the handler
func (p planPoliciesHandler) Name() string {
	return "planpolicies"
}
