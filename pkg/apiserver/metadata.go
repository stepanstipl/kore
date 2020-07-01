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
	"strings"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&metadataHandler{})
}

type metadataHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (c *metadataHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("metadata")
	tags := []string{"metadata"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the metadata webservice")

	c.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllNonValidationErrors(ws.GET("/cloud/{cloud}/regions")).To(c.getCloudRegions).
			Doc("Returns regions for a cloud").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCloudRegions").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve regions for")).
			Returns(http.StatusNotFound, "cloud doesn't exist", nil).
			Returns(http.StatusOK, "A list of all the regions organised by continent", costsv1.ContinentList{}),
	)
	ws.Route(
		withAllNonValidationErrors(ws.GET("/cloud/{cloud}/regions/{region}/zones")).To(c.getCloudRegionZones).
			Doc("Returns supported AZs for a given region of a given cloud provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCloudRegionZones").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve for")).
			Param(ws.PathParameter("region", "The region to retrieve for")).
			Returns(http.StatusNotFound, "cloud or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of supported availability zones", []string{}),
	)
	ws.Route(
		withAllNonValidationErrors(ws.GET("/cloud/{cloud}/regions/{region}/nodetypes")).To(c.getCloudNodeTypes).
			Doc("Returns node types (with prices) for a given region of a given cloud provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCloudNodeTypes").
			Param(ws.PathParameter("cloud", "The cloud provider to retrieve instance types/prices for")).
			Param(ws.PathParameter("region", "The region to retrieve instance types/prices for")).
			Returns(http.StatusNotFound, "cloud or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of instance types with their pricing", costsv1.InstanceTypeList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/k8s/{provider}/regions")).To(c.getKubernetesRegions).
			Doc("Returns regions for a kubernetes provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetKubernetesRegions").
			Param(ws.PathParameter("provider", "The kubernetes provider to retrieve regions for")).
			Returns(http.StatusNotFound, "provider doesn't exist", nil).
			Returns(http.StatusOK, "A list of all the regions organised by continent", costsv1.ContinentList{}),
	)
	ws.Route(
		withAllNonValidationErrors(ws.GET("/k8s/{provider}/regions/{region}/zones")).To(c.getKubernetesRegionZones).
			Doc("Returns supported AZs for a given region of a given kubernetes provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetKubernetesRegionZones").
			Param(ws.PathParameter("provider", "The kubernetes provider to retrieve for")).
			Param(ws.PathParameter("region", "The region to retrieve for")).
			Returns(http.StatusNotFound, "provider or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of supported availability zones", []string{}),
	)
	ws.Route(
		withAllNonValidationErrors(ws.GET("/k8s/{provider}/regions/{region}/instances")).To(c.getKubernetesNodeTypes).
			Doc("Returns node types (with prices) for a given region of a given kubernetes provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetKubernetesNodeTypes").
			Param(ws.PathParameter("provider", "The kubernetes provider to retrieve instance types/prices for")).
			Param(ws.PathParameter("region", "The region to retrieve instance types/prices for")).
			Returns(http.StatusNotFound, "provider or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of instance types with their pricing", costsv1.InstanceTypeList{}),
	)
	ws.Route(
		withAllNonValidationErrors(ws.GET("/k8s/{provider}/regions/{region}/versions")).To(c.getKubernetesVersions).
			Doc("Returns supported Kubernetes versions for a given region of a given kubernetes provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetKubernetesVersions").
			Param(ws.PathParameter("provider", "The kubernetes provider to retrieve for")).
			Param(ws.PathParameter("region", "The region to retrieve for")).
			Returns(http.StatusNotFound, "provider or region doesn't exist", nil).
			Returns(http.StatusOK, "A list of supported kubernetes versions", []string{}),
	)

	return ws, nil
}

func (c metadataHandler) getCloudRegions(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		cloud := strings.ToLower(req.PathParameter("cloud"))
		regions, err := c.Costs().Metadata().Regions(cloud)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, regions)
	})
}
func (c metadataHandler) getCloudRegionZones(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		cloud := strings.ToLower(req.PathParameter("cloud"))
		region := strings.ToLower(req.PathParameter("region"))
		zones, err := c.Costs().Metadata().RegionZones(cloud, region)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, zones)
	})
}
func (c metadataHandler) getCloudNodeTypes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		cloud := strings.ToLower(req.PathParameter("cloud"))
		region := strings.ToLower(req.PathParameter("region"))
		instances, err := c.Costs().Metadata().InstanceTypes(cloud, region)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, instances)
	})
}

func (c metadataHandler) getKubernetesRegions(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		provider := strings.ToUpper(req.PathParameter("provider"))
		cloud, err := c.Costs().Metadata().MapProviderToCloud(provider)
		if err != nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, err)
		}
		regions, err := c.Costs().Metadata().Regions(cloud)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, regions)
	})
}
func (c metadataHandler) getKubernetesRegionZones(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		provider := strings.ToUpper(req.PathParameter("provider"))
		cloud, err := c.Costs().Metadata().MapProviderToCloud(provider)
		if err != nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, err)
		}
		region := strings.ToLower(req.PathParameter("region"))
		zones, err := c.Costs().Metadata().RegionZones(cloud, region)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, zones)
	})
}
func (c metadataHandler) getKubernetesNodeTypes(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		provider := strings.ToUpper(req.PathParameter("provider"))
		cloud, err := c.Costs().Metadata().MapProviderToCloud(provider)
		if err != nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, err)
		}
		region := strings.ToLower(req.PathParameter("region"))
		instances, err := c.Costs().Metadata().InstanceTypes(cloud, region)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, instances)
	})
}
func (c metadataHandler) getKubernetesVersions(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		provider := strings.ToUpper(req.PathParameter("provider"))
		cloud, err := c.Costs().Metadata().MapProviderToCloud(provider)
		if err != nil {
			return resp.WriteHeaderAndEntity(http.StatusNotFound, err)
		}
		region := strings.ToLower(req.PathParameter("region"))
		versions, err := c.Costs().Metadata().KubernetesVersions(cloud, region)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, versions)
	})
}

// Name returns the name of the handler
func (c metadataHandler) Name() string {
	return "metadata"
}
