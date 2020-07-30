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

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&identitiesHandler{})
}

type identitiesHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Name returns the name of the handler
func (u identitiesHandler) Name() string {
	return "identities"
}

// Register is responsible for registering the webserver
func (u *identitiesHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Path("identities")

	log.WithFields(log.Fields{
		"path": path,
	}).Info("registering the identities webservice with container")

	u.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path)

	ws.Route(
		ws.GET("").To(u.findAllIdentities).
			Doc("Returns all the identities for one or more users in kore").
			Operation("ListIdentities").
			Param(ws.QueryParameter("all", "returns all identities managed in kore").DataType("boolean").DefaultValue("false")).
			Param(ws.QueryParameter("type", "set the type of identities to retrieve").DataType("string")).
			Returns(http.StatusOK, "A list of all the identities in the kore", orgv1.IdentityList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{user}")).To(u.findUserIdentities).
			Doc("Find all Identities for a specific user in kore").
			Operation("ListUserIdentity").
			Param(ws.PathParameter("user", "The name of the user you wish to retrieve identities for")).
			Returns(http.StatusOK, "Contains the identities definitions from the kore", orgv1.IdentityList{}).
			Returns(http.StatusNotFound, "User does not exist", nil),
	)

	ws.Route(
		withAllNonValidationErrors(ws.PUT("/{user}/basicauth")).To(u.updateUserBasicAuth).
			Doc("Used to update the basicauth of a local user in kore").
			Operation("UpdateUserBasicauth").
			Param(ws.PathParameter("user", "The name of the user we are updating")).
			Reads(orgv1.UpdateBasicAuthIdentity{}).
			Returns(http.StatusOK, "Contains the identities definitions from the kore", nil).
			Returns(http.StatusNotFound, "User does not exist", nil),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{user}/associate")).To(u.associateUserIdentity).
			Doc("Used to associate an external IDP identity with a user in kore").
			Operation("AssociateIDPIdentity").
			Param(ws.PathParameter("user", "The name of the user you wish to retrieve identities for")).
			Reads(orgv1.UpdateIDPIdentity{}).
			Returns(http.StatusOK, "Contains the identities definitions from the kore", nil).
			Returns(http.StatusNotFound, "User does not exist", nil),
	)

	return ws, nil
}

// updateUserBasicAuth to update the basicauth identity in kore
func (u identitiesHandler) updateUserBasicAuth(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		update := &orgv1.UpdateBasicAuthIdentity{}

		if err := req.ReadEntity(update); err != nil {
			return err
		}

		return u.Users().Identities().UpdateUserBasicAuth(req.Request.Context(), update)
	})
}

// associateUserIdentity to associate an external user identity to the local user
func (u identitiesHandler) associateUserIdentity(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		update := &orgv1.UpdateIDPIdentity{}

		if err := req.ReadEntity(update); err != nil {
			return err
		}

		return u.Users().Identities().AssociateIDPUser(req.Request.Context(), update)
	})
}

// findAllIdentities returns a list of all the identities managed in kore
func (u identitiesHandler) findAllIdentities(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		list, err := u.Users().Identities().List(req.Request.Context(), kore.IdentitiesListOptions{
			IdentityTypes: req.QueryParameters("type"),
		})
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findUserIdentities returns all the identities for a specific user
func (u identitiesHandler) findUserIdentities(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		list, err := u.Users().Identities().List(req.Request.Context(), kore.IdentitiesListOptions{
			User: req.PathParameter("user"),
		})
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}
