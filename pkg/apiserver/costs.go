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
	"time"

	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/validation"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&costHandler{})
}

type costHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (c *costHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("costs")
	tags := []string{"costs"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the costs webservice")

	c.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	// @TODO: admin filters

	ws.Route(
		withAllErrors(ws.POST("")).To(c.postCosts).
			Doc("Persists one or more asset costs").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("PostCosts").
			Reads(costsv1.AssetCostList{}).
			Returns(http.StatusOK, "Costs successfully persisted", nil),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("")).To(c.listCosts).
			Doc("Returns a list of actual costs").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListCosts").
			Param(ws.QueryParameter("team", "Identifier of a team to filter costs for")).
			Param(ws.QueryParameter("asset", "Identifier of an asset to filter costs for")).
			Param(ws.QueryParameter("from", "Start of time range to return costs for")).
			Param(ws.QueryParameter("to", "End of time range to return costs for")).
			Param(ws.QueryParameter("provider", "Cloud provider (e.g. gcp, aws, azure) to return costs for")).
			Param(ws.QueryParameter("account", "Account/project/subscription to return costs for")).
			Returns(http.StatusOK, "A list of costs known to the system, filtered by the above parameters", costsv1.AssetCostList{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("summary/{from}/{to}")).To(c.getCostSummary).
			Doc("Returns a summary of all costs known to kore for the specified time period").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetCostSummary").
			Param(ws.PathParameter("from", "Start of time range to return summary for")).
			Param(ws.PathParameter("to", "End of time range to return summary for")).
			Param(ws.QueryParameter("provider", "Restrict to costs for specified cloud provider (e.g. gcp, aws, azure)")).
			Returns(http.StatusOK, "A summary of costs known to the system", costsv1.OverallCostSummary{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("teamsummary/{teamIdentifier}/{from}/{to}")).To(c.getTeamCostSummary).
			Doc("Returns a summary of all costs known to kore for the specified time period").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetTeamCostSummary").
			Param(ws.PathParameter("teamIdentifier", "Team identifier to retrieve costs for")).
			Param(ws.PathParameter("from", "Start of time range to return summary for")).
			Param(ws.PathParameter("to", "End of time range to return summary for")).
			Param(ws.QueryParameter("provider", "Restrict to costs for specified cloud provider (e.g. gcp, aws, azure)")).
			Returns(http.StatusOK, "A summary of costs known to the system for the team", costsv1.TeamCostSummary{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("assets/{provider}")).To(c.getAssets).
			Doc("Returns details of the assets known to kore which should be monitored for costs by a costs provider").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("GetAssets").
			Param(ws.PathParameter("provider", "Cloud provider (e.g. gcp, aws, azure) to return asset metadata for")).
			Param(ws.QueryParameter("team", "Identifier of a team to filter assets for")).
			Param(ws.QueryParameter("asset", "Identifier of an asset to filter assets for")).
			Param(ws.QueryParameter("with_deleted", "Set to true to include deleted assets")).
			Returns(http.StatusOK, "Metadata describing the assets for the cloud provider in question", costsv1.CostAssetList{}),
	)

	return ws, nil
}

func (c costHandler) postCosts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		costs := &costsv1.AssetCostList{}
		if err := req.ReadEntity(costs); err != nil {
			return err
		}

		err := c.Costs().Assets().StoreAssetCosts(req.Request.Context(), costs)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, nil)
	})
}

func (c costHandler) listCosts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		filters := []persistence.TeamAssetFilterFunc{}
		if req.QueryParameter("from") != "" {
			fromTime, err := time.Parse(time.RFC3339, req.QueryParameter("from"))
			if err != nil {
				return validation.NewError("invalid request").WithFieldErrorf("from", validation.InvalidValue, "cannot parse 'from' time, expected format %s", time.RFC3339)
			}
			filters = append(filters, persistence.TeamAssetFilters.FromTime(fromTime))
		}
		if req.QueryParameter("to") != "" {
			toTime, err := time.Parse(time.RFC3339, req.QueryParameter("to"))
			if err != nil {
				return validation.NewError("invalid request").WithFieldErrorf("to", validation.InvalidValue, "cannot parse 'to' time, expected format %s", time.RFC3339)
			}
			filters = append(filters, persistence.TeamAssetFilters.ToTime(toTime))
		}
		if req.QueryParameter("provider") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithProvider(req.QueryParameter("provider")))
		}
		if req.QueryParameter("team") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithTeam(req.QueryParameter("team")))
		}
		if req.QueryParameter("asset") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithAsset(req.QueryParameter("asset")))
		}
		if req.QueryParameter("provider") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithProvider(req.QueryParameter("provider")))
		}
		if req.QueryParameter("account") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithAccount(req.QueryParameter("account")))
		}
		assets, err := c.Costs().Assets().ListCosts(req.Request.Context(), filters...)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, assets)
	})
}

func (c costHandler) getCostSummary(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		fromTime, err := time.Parse(time.RFC3339, req.PathParameter("from"))
		if err != nil {
			return validation.NewError("invalid request").WithFieldErrorf("from", validation.InvalidValue, "cannot parse 'from' time, expected format %s", time.RFC3339)
		}
		toTime, err := time.Parse(time.RFC3339, req.PathParameter("to"))
		if err != nil {
			return validation.NewError("invalid request").WithFieldErrorf("to", validation.InvalidValue, "cannot parse 'to' time, expected format %s", time.RFC3339)
		}
		filters := []persistence.TeamAssetFilterFunc{}
		if req.QueryParameter("provider") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithProvider(req.QueryParameter("provider")))
		}
		summary, err := c.Costs().Assets().OverallCostsSummary(req.Request.Context(), fromTime, toTime, filters...)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, summary)
	})
}

func (c costHandler) getTeamCostSummary(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		teamIdentifier := req.PathParameter("teamIdentifier")
		fromTime, err := time.Parse(time.RFC3339, req.PathParameter("from"))
		if err != nil {
			return validation.NewError("invalid request").WithFieldErrorf("from", validation.InvalidValue, "cannot parse 'from' time, expected format %s", time.RFC3339)
		}
		toTime, err := time.Parse(time.RFC3339, req.PathParameter("to"))
		if err != nil {
			return validation.NewError("invalid request").WithFieldErrorf("to", validation.InvalidValue, "cannot parse 'to' time, expected format %s", time.RFC3339)
		}
		filters := []persistence.TeamAssetFilterFunc{}
		if req.QueryParameter("provider") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithProvider(req.QueryParameter("provider")))
		}
		summary, err := c.Costs().Assets().TeamCostsSummary(req.Request.Context(), teamIdentifier, fromTime, toTime, filters...)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, summary)
	})
}

func (c costHandler) getAssets(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		filters := []persistence.TeamAssetFilterFunc{
			persistence.TeamAssetFilters.WithProvider(req.PathParameter("provider")),
		}
		if req.QueryParameter("team") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithTeam(req.QueryParameter("team")))
		}
		if req.QueryParameter("asset") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithAsset(req.QueryParameter("asset")))
		}
		if req.QueryParameter("with_deleted") != "" {
			filters = append(filters, persistence.TeamAssetFilters.WithDeleted())
		}
		assets, err := c.Costs().Assets().ListAssets(req.Request.Context(), filters...)
		if err != nil {
			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, assets)
	})
}

// Name returns the name of the handler
func (c costHandler) Name() string {
	return "costs"
}
