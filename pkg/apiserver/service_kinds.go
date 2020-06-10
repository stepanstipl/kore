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

	"github.com/appvia/kore/pkg/apiserver/filters"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&serviceKindsHandler{})
}

type serviceKindsHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *serviceKindsHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("servicekinds")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the servicekinds webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(p.listServiceKinds).
			Doc("Returns all the available service kinds").
			Operation("ListServiceKinds").
			Param(ws.QueryParameter("platform", "Filters service kinds for a specific service platform")).
			Param(ws.QueryParameter("enabled", "Filters service kinds for enabled/disabled status")).
			Returns(http.StatusOK, "A list of service kinds", servicesv1.ServiceKindList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.getServiceKind).
			Doc("Returns a specific service kind").
			Operation("GetServiceKind").
			Param(ws.PathParameter("name", "The name of the service kind you wish to retrieve")).
			Returns(http.StatusNotFound, "the service kind with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service kind definition", servicesv1.ServiceKind{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updateServiceKind).
			Filter(filters.Admin).
			Filter(p.readOnlyServiceKindFilter).
			Doc("Creates or updates a service kind").
			Operation("UpdateServiceKind").
			Param(ws.PathParameter("name", "The name of the service kind you wish to create or update")).
			Reads(servicesv1.ServiceKind{}, "The specification for the service kind you are creating or updating").
			Returns(http.StatusOK, "Contains the service kind definition", servicesv1.ServiceKind{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deleteServiceKind).
			Filter(filters.Admin).
			Doc("Deletes a service kind").
			Operation("DeleteServiceKind").
			Param(ws.PathParameter("name", "The name of the service kind you wish to delete")).
			Returns(http.StatusNotFound, "the service kind with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service kind definition", servicesv1.ServiceKind{}),
	)

	return ws, nil
}

func (p serviceKindsHandler) readOnlyServiceKindFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		serviceKind, err := p.ServiceKinds().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if serviceKind != nil && serviceKind.Annotations[kore.AnnotationReadOnly] == "true" {
			resp.WriteHeader(http.StatusForbidden)
			return nil
		}

		// @step: continue with the chain
		chain.ProcessFilter(req, resp)
		return nil
	})
}

// getServiceKind returns a specific service kind
func (p serviceKindsHandler) getServiceKind(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		kind, err := p.ServiceKinds().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, kind)
	})
}

// listServiceKinds returns all service kinds in the kore
func (p serviceKindsHandler) listServiceKinds(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		var filters []func(servicesv1.ServiceKind) bool
		if platform := req.QueryParameter("platform"); platform != "" {
			filters = append(filters, func(s servicesv1.ServiceKind) bool {
				return s.Labels[kore.Label("platform")] == platform
			})
		}
		if enabled := req.QueryParameter("enabled"); enabled != "" {
			filters = append(filters, func(s servicesv1.ServiceKind) bool {
				return s.Spec.Enabled == (enabled == "true")
			})
		}

		res, err := p.ServiceKinds().List(req.Request.Context(), filters...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, res)
	})
}

// updateServiceKind is used to update or create a service kind in the kore
func (p serviceKindsHandler) updateServiceKind(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		kind := &servicesv1.ServiceKind{}
		if err := req.ReadEntity(kind); err != nil {
			return err
		}
		kind.Name = name

		existing, err := p.ServiceKinds().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if existing == nil {
			writeError(req, resp, fmt.Errorf("creating new service kinds are not allowed"), http.StatusMethodNotAllowed)
			return nil
		}

		if err := p.ServiceKinds().Update(req.Request.Context(), kind); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, kind)
	})
}

// deleteServiceKind is used to update or create a service kind in the kore
func (p serviceKindsHandler) deleteServiceKind(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		writeError(req, resp, fmt.Errorf("deleting service kinds are not allowed, please set enabled to false instead"), http.StatusMethodNotAllowed)
		return nil
	})
}

// Name returns the name of the handler
func (p serviceKindsHandler) Name() string {
	return "servicekinds"
}

// Enabled returns true if the services feature gate is enabled
func (p serviceKindsHandler) Enabled() bool {
	return p.Config().IsFeatureGateEnabled(kore.FeatureGateServices)
}
