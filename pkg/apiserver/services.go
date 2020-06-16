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

	"github.com/appvia/kore/pkg/kore/authentication"

	"github.com/appvia/kore/pkg/kore"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	restful "github.com/emicklei/go-restful"
)

func (u teamHandler) readonlyServiceFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		service, err := u.Teams().Team(team).Services().Get(req.Request.Context(), name)
		if err != nil && err != kore.ErrNotFound {
			return err
		}

		if service != nil && service.Annotations[kore.AnnotationReadOnly] == "true" {
			resp.WriteHeader(http.StatusForbidden)
			return nil
		}

		// @step: continue with the chain
		chain.ProcessFilter(req, resp)
		return nil
	})
}

// listServices returns all the services from a team
func (u teamHandler) listServices(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		var list *servicesv1.ServiceList
		var err error

		user := authentication.MustGetIdentity(req.Request.Context())
		if user.IsGlobalAdmin() {
			list, err = u.Teams().Team(team).Services().List(req.Request.Context())
		} else {
			list, err = u.Teams().Team(team).Services().ListFiltered(req.Request.Context(), func(service servicesv1.Service) bool {
				return service.Annotations[kore.AnnotationSystem] != "true"
			})
		}
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getService returns a service from a team
func (u teamHandler) getService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		service, err := u.Teams().Team(team).Services().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, service)
	})
}

// updateService is responsible for creating or updating a service
func (u teamHandler) updateService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		service := &servicesv1.Service{}
		if err := req.ReadEntity(service); err != nil {
			return err
		}

		if service.Annotations[kore.AnnotationReadOnly] != "" {
			writeError(req, resp, fmt.Errorf("setting %q annotation is not allowed", kore.AnnotationReadOnly), http.StatusForbidden)
			return nil
		}

		if err := u.Teams().Team(team).Services().Update(req.Request.Context(), service); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, service)
	})
}

// deleteService is responsible for deleting a service from a team
func (u teamHandler) deleteService(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Services().Delete(ctx, name, parseDeleteOpts(req)...)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
