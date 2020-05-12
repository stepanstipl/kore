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
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	RegisterHandler(&securityScansHandler{})
}

type securityScansHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (s *securityScansHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("securityscans")
	tags := []string{"security"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the security scans webservice")

	s.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(s.listSecurityScans).
			Doc("Used to return a list of security scan results").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListSecurityScans").
			Param(ws.QueryParameter("latestOnly", "Set to false to retrieve full history").DefaultValue("true").DataType("boolean")).
			Returns(http.StatusOK, "A collection of security scans", securityv1.SecurityScanResultList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("overview")).To(s.getSecurityOverview).
			Doc("Used to return a summary of the security overview for the entire Kore estate").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetSecurityOverview").
			Returns(http.StatusOK, "A report of the security posture of Kore", securityv1.SecurityOverview{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("{group}/{version}/{kind}/{namespace}/{name}")).To(s.getSecurityScanForResource).
			Doc("Used to return latest security scan for specific object in the system").
			Metadata(restfulspec.KeyOpenAPITags, tags).
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
		withAllErrors(ws.PUT("scans/{group}/{version}/{kind}/{namespace}/{name}")).To(s.storeSecurityScanForResource).
			Doc("Used to persist a new security scan result for specific object in the system").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("StoreSecurityScanForResource").
			Reads(securityv1.SecurityScanResult{}, "The result of a security scan").
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("name", "Is the name of the resource")).
			Returns(http.StatusOK, "Latest security scan for this resource", securityv1.SecurityScanResult{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("scans/{group}/{version}/{kind}/{namespace}/{name}/history")).
			To(s.listSecurityScansForResource).
			Doc("Used to return the history of security scans for specific object in the system").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListSecurityScansForResource").
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("name", "Is the name of the resource")).
			Returns(http.StatusOK, "Latest security scan for this resource", securityv1.SecurityScanResultList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("{id}")).To(s.getSecurityScan).
			Doc("Used to return specific security scan by ID").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetSecurityScan").
			Param(ws.PathParameter("id", "Is the ID of the scan").DataType("integer")).
			Returns(http.StatusNotFound, "No current security scan exists for the ID", nil).
			Returns(http.StatusOK, "Security scan", securityv1.SecurityScanResult{}),
	)

	return ws, nil
}

func (s *securityScansHandler) listSecurityScans(req *restful.Request, resp *restful.Response) {
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

func (s *securityScansHandler) getSecurityOverview(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		overview, err := s.Security().GetOverview(req.Request.Context())
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, overview)
	})
}

func (s *securityScansHandler) getSecurityScanForResource(req *restful.Request, resp *restful.Response) {
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

func (s *securityScansHandler) storeSecurityScanForResource(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {

		scanResult := &securityv1.SecurityScanResult{}
		if err := req.ReadEntity(scanResult); err != nil {
			return err
		}

		return fmt.Errorf("Not implemented")
	})
}

func (s *securityScansHandler) listSecurityScansForResource(req *restful.Request, resp *restful.Response) {
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

func (s *securityScansHandler) getSecurityScan(req *restful.Request, resp *restful.Response) {
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
func (s securityScansHandler) Name() string {
	return "securityscans"
}
