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

	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&securityRulesHandler{})
}

type securityRulesHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (s *securityRulesHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("securityrules")
	tags := []string{"security"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the security rules webservice")

	s.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(s.listSecurityRules).
			Doc("Used to return a list of all the security rules in the system").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListSecurityRules").
			Returns(http.StatusOK, "A collection of security rules", securityv1.SecurityRuleList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("{code}")).To(s.getSecurityRule).
			Doc("Used to return details of a specific security rule within the system").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetSecurityRule").
			Param(ws.PathParameter("code", "Is the unique code for the security rule")).
			Returns(http.StatusNotFound, "No security rule exists for the code", nil).
			Returns(http.StatusOK, "A security rule", securityv1.SecurityRule{}),
	)

	return ws, nil
}

func (s *securityRulesHandler) listSecurityRules(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		list, err := s.Security().ListRules(req.Request.Context())
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

func (s *securityRulesHandler) getSecurityRule(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		rule, err := s.Security().GetRule(req.Request.Context(), req.PathParameter("code"))
		if err != nil {
			return err
		}
		if rule == nil {
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, rule)
	})
}

// Name returns the name of the handler
func (s securityRulesHandler) Name() string {
	return "securityrules"
}
