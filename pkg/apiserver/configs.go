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

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&configsHandler{})
}

type configsHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (v *configsHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Path("configs")

	log.WithFields(log.Fields{
		"path": path,
	}).Info("registering the config webservice with container")

	v.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path)

	ws.Route(
		ws.GET("").To(v.listConfig).
			Doc("Returns all the configs in the kore").
			Operation("ListConfig").
			Returns(http.StatusOK, "A list of all the config values in the kore", configv1.ConfigList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{config}").To(v.getConfig)).
			Doc("Return information related to the specific config name in kore").
			Operation("GetConfig").
			Param(ws.PathParameter("config", "The name of the config you wish to retrieve")).
			Returns(http.StatusOK, "A list of all the config in the kore", configv1.Config{}).
			Returns(http.StatusNotFound, "config does not exist", nil),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{config}")).To(v.updateConfig).
			Doc("Used to create or update a config in the kore").
			Operation("UpdateConfig").
			Param(ws.PathParameter("config", "The name of the config you are updating or creating in the kore")).
			Reads(configv1.Config{}, "The specification for a config in the kore").
			Returns(http.StatusOK, "Contains the config definition from the kore", configv1.Config{}).
			Returns(http.StatusNotFound, "config does not exist", nil),
	)

	ws.Route(
		ws.DELETE("/{config}").To(v.deleteConfig).
			Doc("Used to delete a config from the kore").
			Operation("RemoveConfig").
			Param(ws.PathParameter("config", "The name of the config you are deleting from the kore")).
			Returns(http.StatusOK, "Contains the former config definition from the kore", configv1.Config{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

func (v configsHandler) listConfig(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		configs, err := v.Configs().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, configs)
	})

}

func (v configsHandler) getConfig(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		config, err := v.Configs().Get(req.Request.Context(), req.PathParameter("config"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, config)
	})
}

func (v configsHandler) updateConfig(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		config := &configv1.Config{}
		if err := req.ReadEntity(config); err != nil {
			return err
		}

		config, err := v.Configs().Update(req.Request.Context(), config)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, config)
	})
}

func (v configsHandler) deleteConfig(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("config")

		config, err := v.Configs().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, config)
	})
}

func (v configsHandler) Name() string {
	return "config"
}
