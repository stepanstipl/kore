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

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&idpHandler{})
}

type idpHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (id idpHandler) Name() string {
	return "idp"
}

// Register is responsible for handling the registration
func (id *idpHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.Info("registering the idp webservice")
	id.Interface = i
	path := builder.Add("idp")

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the idp webservice")

	// Types of IDP providers that can be configured
	ws.Route(
		ws.GET("/types").To(id.getTypes).
			Doc("Returns a list of all the possible identity providers supported in the kore").
			Operation("ListIDPTypes").
			Returns(http.StatusOK, "A list of all the possible identity provider types", []corev1.IDPConfig{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// The default IDP provider for identity in the kore
	ws.Route(
		ws.GET("/default").To(id.getDefaultIDP).
			Doc("Returns the default identity provider configured in the kore").
			Operation("GetDefaultIDP").
			Returns(http.StatusOK, "The default configured identity provider", corev1.IDP{}).
			Returns(http.StatusNotFound, "Indicate the class was not found in the kore", nil).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// All configured IDP providers
	ws.Route(
		ws.GET("/configured/").To(id.findIDPs).
			Doc("Returns a list of all the configured identity providers in the kore").
			Operation("ListIDPs").
			Returns(http.StatusOK, "A list of all the configured identity providers", []corev1.IDP{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/configured/{name}").To(id.getIDP).
			Doc("Returns the definition for a specific identity provider").
			Operation("GetIDP").
			Param(ws.PathParameter("name", "The name of the configured IDP provider to retrieve")).
			Returns(http.StatusOK, "the specified identity provider", corev1.IDP{}).
			Returns(http.StatusNotFound, "Indicate the class was not found in the kore", nil).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/configured/{name}").To(id.putIDP).
			Doc("Returns the definition for a specific ID provider").
			Operation("UpdateIDP").
			Param(ws.PathParameter("name", "The name of the configured IDP provider to update")).
			Reads(corev1.IDP{}, "The definition for the ID provider").
			Returns(http.StatusOK, "A list of all the IDPs in the kore", corev1.IDP{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// CLients confgured by name in the kore
	ws.Route(
		ws.PUT("/clients/{name}").To(id.putIDPClient).
			Doc("Updates the definition for a specific idp client").
			Operation("UpdateIDPClient").
			Param(ws.PathParameter("name", "The name of the IDP client provider to update")).
			Reads(corev1.IDPClient{}, "The definition for the idp client").
			Returns(http.StatusOK, "The configured client in the kore", corev1.IDPClient{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// getTypes
func (id idpHandler) getTypes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		c := id.IDP().ConfigTypes(req.Request.Context())
		return resp.WriteHeaderAndEntity(http.StatusOK, c)
	})
}

// findIDPs returns all the providers in the kore
func (id idpHandler) findIDPs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		idp, err := id.IDP().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, idp)
	})
}

// getIDP returns a specific provider from the kore
func (id idpHandler) getIDP(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		idp, err := id.IDP().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, idp)
	})
}

// getDefaultIDP returns a specific provider from the kore
func (id idpHandler) getDefaultIDP(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		idp, err := id.IDP().Default(req.Request.Context())
		if err == kore.ErrNotFound {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, idp)
		}
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, idp)
	})
}

// putIDP updates the auth provider
func (id idpHandler) putIDP(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		idp := &corev1.IDP{}
		if err := req.ReadEntity(idp); err != nil {
			return err
		}
		idp.SetName(name)

		if err := id.IDP().Update(req.Request.Context(), idp); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, idp)
	})
}

// putIDPClient updates the IDP client
func (id idpHandler) putIDPClient(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		c := &corev1.IDPClient{}
		if err := req.ReadEntity(c); err != nil {
			return err
		}
		c.SetName(name)

		if err := id.IDP().UpdateClient(req.Request.Context(), c); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, c)
	})
}
