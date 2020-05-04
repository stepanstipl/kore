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
	"strconv"

	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/validation"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	RegisterHandler(&securityHandler{})
}

type securityHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (s *securityHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("security")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the security webservice")

	s.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("rules")).To(s.listSecurityRules).
			Doc("Used to return a list of all the security rules in the system").
			Operation("ListSecurityRules").
			Returns(http.StatusOK, "A collection of security rules", securityv1.SecurityRuleList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("rules/{code}")).To(s.getSecurityRule).
			Doc("Used to return details of a specific security rule within the system").
			Operation("GetSecurityRule").
			Param(ws.PathParameter("code", "Is the unique code for the security rule")).
			Returns(http.StatusNotFound, "No security rule exists for the code", nil).
			Returns(http.StatusOK, "A security rule", securityv1.SecurityRule{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("scans")).To(s.listSecurityScans).
			Doc("Used to return security scans for any object in the system").
			Operation("ListSecurityScans").
			Param(ws.QueryParameter("latestOnly", "Set to false to retrieve full history").DefaultValue("true").DataType("boolean")).
			Returns(http.StatusOK, "A collection of security scans", securityv1.SecurityScanResultList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("scans/{group}/{version}/{kind}/{namespace}/{name}")).To(s.getSecurityScanForResource).
			Doc("Used to return latest security scan for specific object in the system").
			Operation("GetSecurityScanForResource").
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("name", "Is the name of the resource")).
			Returns(http.StatusNotFound, "No current security scan exists for the resource", nil).
			Returns(http.StatusOK, "Latest security scan for this resource", securityv1.SecurityScanResult{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("scans/{group}/{version}/{kind}/{namespace}/{name}/history")).
			To(s.listSecurityScansForResource).
			Doc("Used to return the history of security scans for specific object in the system").
			Operation("ListSecurityScansForResource").
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("name", "Is the name of the resource")).
			Returns(http.StatusOK, "Latest security scan for this resource", securityv1.SecurityScanResultList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("scans/{id}")).To(s.getSecurityScan).
			Doc("Used to return specific security scan by ID").
			Operation("GetSecurityScan").
			Param(ws.PathParameter("id", "Is the ID of the scan").DataType("integer")).
			Returns(http.StatusNotFound, "No current security scan exists for the ID", nil).
			Returns(http.StatusOK, "Security scan", securityv1.SecurityScanResult{}),
	)

	return ws, nil
}

func (s *securityHandler) listSecurityRules(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		list, err := s.Security().ListRules(req.Request.Context())
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

func (s *securityHandler) getSecurityRule(req *restful.Request, resp *restful.Response) {
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

func (s *securityHandler) listSecurityScans(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		latestOnly, err := strconv.ParseBool(req.QueryParameter("latestOnly"))
		if err != nil {
			return validation.NewError("Invalid request").
				WithFieldError("latestOnly", validation.InvalidType, "should be a boolean")
		}
		list, err := s.Security().ListScans(req.Request.Context(), latestOnly)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

func (s *securityHandler) getSecurityScanForResource(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		scan, err := s.Security().GetCurrentScanForResource(
			req.Request.Context(),
			metav1.TypeMeta{
				APIVersion: fmt.Sprintf("%s/%s", req.PathParameter("group"), req.PathParameter("version")),
				Kind:       req.PathParameter("kind"),
			},
			metav1.ObjectMeta{
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("name"),
			},
		)
		if err != nil {
			return err
		}
		if scan == nil {
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, scan)
	})
}

func (s *securityHandler) listSecurityScansForResource(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		list, err := s.Security().ScanHistoryForResource(
			req.Request.Context(),
			metav1.TypeMeta{
				APIVersion: fmt.Sprintf("%s/%s", req.PathParameter("group"), req.PathParameter("version")),
				Kind:       req.PathParameter("kind"),
			},
			metav1.ObjectMeta{
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("name"),
			},
		)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

func (s *securityHandler) getSecurityScan(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		id, err := strconv.ParseUint(req.PathParameter("id"), 10, 64)
		if err != nil {
			return err
		}
		scan, err := s.Security().GetScan(req.Request.Context(), id)
		if err != nil {
			return err
		}
		if scan == nil {
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, scan)
	})
}

// Name returns the name of the handler
func (s securityHandler) Name() string {
	return "security"
}
