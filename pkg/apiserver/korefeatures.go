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
	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&koreFeaturesHandler{})
}

type koreFeaturesHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (f *koreFeaturesHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("korefeatures")
	tags := []string{"korefeatures"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the korefeatures webservice")

	f.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(f.listFeatures).
			Doc("Returns a list of features").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListFeatures").
			Filter(filters.Admin).
			Returns(http.StatusOK, "A list of all features known to the system", configv1.KoreFeatureList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{name}")).To(f.getFeature).
			Doc("Returns a specific feature").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetFeature").
			Filter(filters.Admin).
			Param(ws.PathParameter("name", "The name of the feature you wish to retrieve")).
			Returns(http.StatusNotFound, "the feature with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the feature definition", configv1.KoreFeature{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{name}")).To(f.updateFeature).
			Doc("Used to create or update a feature").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("UpdateFeature").
			Filter(filters.Admin).
			Param(ws.PathParameter("name", "The name of the feature you wish to create or update")).
			Reads(configv1.KoreFeature{}, "The specification for the feature you are creating or updating").
			Returns(http.StatusOK, "Contains the feature definition", configv1.KoreFeature{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/{name}")).To(f.deleteFeature).
			Doc("Used to delete a feature").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("RemoveFeature").
			Filter(filters.Admin).
			Param(ws.PathParameter("name", "The name of the feature you wish to delete")).
			Returns(http.StatusNotFound, "the feature with the given name doesn't exist", nil).
			Returns(http.StatusOK, "Contains the feature definition", configv1.KoreFeature{}),
	)

	return ws, nil
}

func (f koreFeaturesHandler) listFeatures(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		features, err := f.Features().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, features)
	})
}

func (f koreFeaturesHandler) getFeature(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		feature, err := f.Features().Get(req.Request.Context(), req.PathParameter("name"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, feature)
	})
}

func (f koreFeaturesHandler) updateFeature(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		feature := &configv1.KoreFeature{}
		if err := req.ReadEntity(feature); err != nil {
			return err
		}
		feature.Name = name

		feat, err := f.Features().Update(req.Request.Context(), feature)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, feat)
	})
}

func (f koreFeaturesHandler) deleteFeature(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		name := req.PathParameter("name")

		feature, err := f.Features().Delete(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, feature)
	})
}

// Name returns the name of the handler
func (f koreFeaturesHandler) Name() string {
	return "korefeatures"
}
