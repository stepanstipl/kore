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

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"

	restful "github.com/emicklei/go-restful"
)

// listServiceCredentials returns all the service credentials from a team
func (u teamHandler) listServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).ServiceCredentials().List(req.Request.Context())
		if err != nil {
			return err
		}

		var filters []func(servicesv1.ServiceCredentials) bool

		cluster := req.QueryParameter("cluster")
		if cluster != "" {
			filters = append(filters, func(s servicesv1.ServiceCredentials) bool { return s.Spec.Cluster.Name == cluster })
		}

		service := req.QueryParameter("service")
		if service != "" {
			filters = append(filters, func(s servicesv1.ServiceCredentials) bool { return s.Spec.Service.Name == service })
		}

		if len(filters) > 0 {
			filtered := []servicesv1.ServiceCredentials{}
			for _, item := range list.Items {
				include := true
				for _, filter := range filters {
					include = include && filter(item)
				}
				if include {
					filtered = append(filtered, item)
				}
			}
			list.Items = filtered
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// getServiceCredentials returns service credentials from a team
func (u teamHandler) getServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).ServiceCredentials().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// updateServiceCredentials is responsible for creating or updating service credentials
func (u teamHandler) updateServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		serviceCreds := &servicesv1.ServiceCredentials{}
		if err := req.ReadEntity(serviceCreds); err != nil {
			return err
		}

		if err := u.Teams().Team(team).ServiceCredentials().Update(req.Request.Context(), serviceCreds); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, serviceCreds)
	})
}

// deleteServiceCredentials is responsible for deleting service credentials from a team
func (u teamHandler) deleteServiceCredentials(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).ServiceCredentials().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
