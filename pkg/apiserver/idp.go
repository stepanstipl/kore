/*
 * Copyright (C) 2019 Appvia Ltd. <lewis.marshall@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"net/http"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&idpHandler{})
}

type idpHandler struct {
	hub.Interface
	// DefaultHandlder implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (id idpHandler) Name() string {
	return "idp"
}

// Register is responsible for handling the registration
func (id *idpHandler) Register(i hub.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	log.Info("registering the idp webservice")
	id.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(builder.Path("idp"))

	// Types of IDP providers that can be configured
	ws.Route(
		ws.GET("/types").To(id.getTypes).
			Doc("Returns a list of all the possible identity providers supported in the hub").
			Returns(http.StatusOK, "A list of all the possible identity providers types", []corev1.IDPConfig{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// The default IDP provider for identity in the hub
	ws.Route(
		ws.GET("/default").To(id.getDefaultIDP).
			Doc("Returns the default identity provider configured in the hub").
			Returns(http.StatusOK, "The default configured identity provider", corev1.IDP{}).
			Returns(http.StatusNotFound, "Indicate the class was not found in the hub", nil).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// All configured IDP providers
	ws.Route(
		ws.GET("/configured/").To(id.findIDPs).
			Doc("Returns a list of all the configured identity providers in the hub").
			Returns(http.StatusOK, "A list of all the configured identity providers", []corev1.IDP{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/configured/{name}").To(id.getIDP).
			Doc("Returns the definition for a specific identity provider").
			Param(ws.PathParameter("name", "The name of the configured IDP provider to retrieve")).
			Returns(http.StatusOK, "the specified identity provider", corev1.IDP{}).
			Returns(http.StatusNotFound, "Indicate the class was not found in the hub", nil).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/configured/{name}").To(id.putIDP).
			Doc("Returns the definition for a specific ID provider").
			Param(ws.PathParameter("name", "The name of the configured IDP provider to update")).
			Reads(corev1.IDP{}, "The definition for the ID provider").
			Returns(http.StatusOK, "A list of all the IDPs in the hub", corev1.IDP{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// CLients confgured by name in the hub
	ws.Route(
		ws.PUT("/clients/{name}").To(id.putIDPClient).
			Doc("Returns the definition for a specific idp client").
			Param(ws.PathParameter("name", "The name of the IDP client provider to update")).
			Reads(corev1.IDPClient{}, "The definition for the idp client").
			Returns(http.StatusOK, "The configured client in the hub", corev1.IDPClient{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// getTypes
func (id idpHandler) getTypes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		c, err := id.IDP().ConfigTypes(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, c)
	})
}

// findIDPs returns all the providers in the hub
func (id idpHandler) findIDPs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		idp, err := id.IDP().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, idp)
	})
}

// getIDP returns a specific provider from the hub
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

// getDefaultIDP returns a specific provider from the hub
func (id idpHandler) getDefaultIDP(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		idp, err := id.IDP().Default(req.Request.Context())
		if err == hub.ErrNotFound {
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
