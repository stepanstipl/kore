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

	accountsv1beta1 "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&accountsHandler{})
}

type accountsHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (p *accountsHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("accountmanagements")

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the accounts webservice")

	p.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(p.findAccounts).
			Doc("Returns all the accounts available to initialized in the kore").
			Operation("ListAccounts").
			Param(ws.QueryParameter("kind", "Returns all accounts for a specific resource type")).
			Returns(http.StatusOK, "A list of all the accounts", accountsv1beta1.AccountManagementList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(p.findAccount).
			Doc("Returns a specific account account from the kore").
			Operation("GetAccount").
			Param(ws.PathParameter("name", "The name of the account you wish to retrieve")).
			Returns(http.StatusNotFound, "the account with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the account definition", accountsv1beta1.AccountManagement{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(p.updateAccount).
			Doc("Used to create or update a account in the kore").
			Operation("UpdateAccount").
			Param(ws.PathParameter("name", "The name of the account you wish to create or update")).
			Reads(accountsv1beta1.AccountManagement{}, "The specification for the account you are creating or updating").
			Returns(http.StatusOK, "Contains the account definition", accountsv1beta1.AccountManagement{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(p.deleteAccount).
			Doc("Used to delete a account from the kore").
			Operation("RemoveAccount").
			Param(ws.PathParameter("name", "The name of the account you wish to delete")).
			Returns(http.StatusNotFound, "the account with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the account definition", accountsv1beta1.AccountManagement{}),
	)

	return ws, nil
}

// findAccount returns a specific account
func (p accountsHandler) findAccount(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		account, err := p.Accounts().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, account)
	})
}

// findAccounts returns all accounts in the kore
func (p accountsHandler) findAccounts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		accounts, err := p.Accounts().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, accounts)
	})
}

// updateAccount is used to update or create a account in the kore
func (p accountsHandler) updateAccount(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		account := &accountsv1beta1.AccountManagement{}
		if err := req.ReadEntity(account); err != nil {
			return err
		}
		account.Name = name

		if err := p.Accounts().Update(req.Request.Context(), account); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, account)
	})
}

// deleteAccount is used to update or create a account in the kore
func (p accountsHandler) deleteAccount(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		account, err := p.Accounts().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, account)
	})
}

// Name returns the name of the handler
func (p accountsHandler) Name() string {
	return "accounts"
}
