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

	"github.com/appvia/kore/pkg/apiserver/params"

	"github.com/appvia/kore/pkg/apiserver/filters"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&serviceProvidersHandler{})
}

type serviceProvidersHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

func (p *serviceProvidersHandler) readOnlyServiceProviderFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		serviceProvider, err := p.ServiceProviders().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if serviceProvider != nil && serviceProvider.Annotations[kore.AnnotationReadOnly] == "true" {
			resp.WriteHeader(http.StatusForbidden)
			return nil
		}

		// @step: continue with the chain
		chain.ProcessFilter(req, resp)
		return nil
	})
}

// Register is called by the api server on registration
func (p *serviceProvidersHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("serviceproviders")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the serviceproviders webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(p.findServiceProviders).
			Doc("Returns all the available service providers").
			Operation("ListServiceProviders").
			Param(ws.QueryParameter("kind", "Filters service providers for a specific kind")).
			Returns(http.StatusOK, "A list of service providers", servicesv1.ServiceProviderList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.findServiceProvider).
			Doc("Returns a specific service provider").
			Operation("GetServiceProvider").
			Param(ws.PathParameter("name", "The name of the service provider you wish to retrieve")).
			Returns(http.StatusNotFound, "the service provider with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service provider definition", servicesv1.ServiceProvider{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updateServiceProvider).
			Filter(filters.Admin).
			Filter(p.readOnlyServiceProviderFilter).
			Doc("Creates or updates a service provider").
			Operation("UpdateServiceProvider").
			Param(ws.PathParameter("name", "The name of the service provider you wish to create or update")).
			Reads(servicesv1.ServiceProvider{}, "The specification for the service provider you are creating or updating").
			Returns(http.StatusOK, "Contains the service provider definition", servicesv1.ServiceProvider{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deleteServiceProvider).
			Filter(filters.Admin).
			Filter(p.readOnlyServiceProviderFilter).
			Doc("Deletes a service provider").
			Operation("DeleteServiceProvider").
			Param(ws.PathParameter("name", "The name of the service provider you wish to delete")).
			Param(params.DeleteCascade()).
			Returns(http.StatusNotFound, "the service provider with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the service provider definition", servicesv1.ServiceProvider{}),
	)

	return ws, nil
}

// findServiceProvider returns a specific service provider
func (p serviceProvidersHandler) findServiceProvider(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		provider, err := p.ServiceProviders().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, provider)
	})
}

// findServiceProviders returns all service providers in the kore
func (p serviceProvidersHandler) findServiceProviders(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		res, err := p.ServiceProviders().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, res)
	})
}

// updateServiceProvider is used to update or create a service provider in the kore
func (p serviceProvidersHandler) updateServiceProvider(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		provider := &servicesv1.ServiceProvider{}
		if err := req.ReadEntity(provider); err != nil {
			return err
		}
		provider.Name = name

		if provider.Annotations[kore.AnnotationReadOnly] != "" {
			writeError(req, resp, fmt.Errorf("setting %q annotation is not allowed", kore.AnnotationReadOnly), http.StatusForbidden)
			return nil
		}

		if err := p.ServiceProviders().Update(req.Request.Context(), provider); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, provider)
	})
}

// deleteServiceProvider is used to update or create a service provider in the kore
func (p serviceProvidersHandler) deleteServiceProvider(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		provider, err := p.ServiceProviders().Delete(req.Request.Context(), name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, provider)
	})
}

// Name returns the name of the handler
func (p serviceProvidersHandler) Name() string {
	return "serviceproviders"
}

// Enabled returns true if the services feature gate is enabled
func (p serviceProvidersHandler) Enabled() bool {
	return p.Config().IsFeatureGateEnabled(kore.FeatureGateServices)
}
