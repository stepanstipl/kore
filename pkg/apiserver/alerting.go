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
	"strconv"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	monitoring "github.com/appvia/kore/pkg/apis/monitoring/v1beta1"
	"github.com/appvia/kore/pkg/apiserver/filters"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&alertsHandler{})
}

type alertsHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server on registration
func (a *alertsHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	path := builder.Add("monitoring")
	tags := []string{"monitoring"}

	log.WithFields(log.Fields{
		"path": path.Base(),
	}).Info("registering the alerts webservice")

	a.Interface = i

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path.Base())

	ws.Route(
		withAllErrors(ws.GET("/rules")).To(a.findAllRules).
			Filter(filters.Admin).
			Doc("Returns all available rules currently in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("ListRules").
			Returns(http.StatusOK, "Listing of the rules in kore", monitoring.RuleList{}),
	)

	ws.Route(
		withAllErrors(ws.GET("/rules/{group}/{version}/{kind}/{namespace}/{resource}")).To(a.findResourceRules).
			Filter(filters.Admin).
			Doc("Get all the on a specific resource in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.QueryParameter("source", "The producer of the alerting rule")).
			Operation("GetRules").
			Returns(http.StatusOK, "The rule has been deleted", monitoring.Rule{}),
	)

	ws.Route(
		withAllErrors(ws.GET("/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}")).To(a.findRule).
			Filter(filters.Admin).
			Doc("Returns the definition of a rule in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.PathParameter("name", "Is the name of the alerting rule")).
			Operation("GetRule").
			Returns(http.StatusOK, "The definition of the monitoring rule", monitoring.Rule{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}")).To(a.updateRule).
			Filter(filters.Admin).
			Reads(monitoring.Rule{}, "The specification for a rule in the kore").
			Doc("Updates or creates a rule in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.PathParameter("name", "Is the name of the alerting rule")).
			Operation("UpdateRule").
			Returns(http.StatusOK, "The rule has been deleted", monitoring.Rule{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/rules/{group}/{version}/{kind}/{namespace}/{resource}")).To(a.deleteResourceRules).
			Filter(filters.Admin).
			Doc("Deletes all the rules on a specific resource").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.QueryParameter("source", "The producer of the alerting rule")).
			Operation("DeleteResourceRules").
			Returns(http.StatusOK, "The rules have been deleted", nil),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}")).To(a.deleteRule).
			Filter(filters.Admin).
			Reads(monitoring.Rule{}, "The specification for a rule in the kore").
			Doc("Deletes a rule and all history from a resource").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.PathParameter("name", "Is the name of the alerting rule")).
			Operation("DeleteRule").
			Returns(http.StatusOK, "The rule has been deleted", monitoring.Rule{}),
	)

	ws.Route(
		withAllErrors(ws.DELETE("/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}/purge")).To(a.purgeRule).
			Filter(filters.Admin).
			Doc("Purges the history of alerts or a given rule").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.PathParameter("name", "Is the name of the alerting rule")).
			Param(ws.QueryParameter("duration", "The duration to keep i.e. 1h 24d")).
			Operation("PurgeRuleAlerts").
			Returns(http.StatusOK, "The history has been purged", nil),
	)

	ws.Route(
		withAllErrors(ws.GET("/alerts")).To(a.findAllAlerts).
			Filter(filters.Admin).
			Doc("Returns all available alerts currently in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.QueryParameter("status", "The alert to filter the results by")).
			Param(ws.QueryParameter("history", "The number of historical records to retrieve")).
			Param(ws.QueryParameter("label", "A label to filter the alert by")).
			Param(ws.QueryParameter("latest", "Indicates to we only want the latest alert status")).
			Operation("ListAlerts").
			Returns(http.StatusOK, "Listing of the alert statuses in kore", monitoring.AlertList{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/alerts/status")).To(a.updateAlert).
			Filter(filters.Admin).
			Doc("Used to store an alert in kore").
			Reads(monitoring.Alert{}, "The definition of the alert being created or updated").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Operation("UpdateAlert").
			Returns(http.StatusOK, "The alert has been successfully stored", nil),
	)

	ws.Route(
		withAllErrors(ws.GET("/alerts/resource/{group}/{version}/{kind}/{namespace}/{resource}")).To(a.findAlertsOnResource).
			Filter(filters.Admin).
			Doc("Used to retrieve the alerts on a resource").
			Reads(monitoring.Alert{}, "The definition of the alert being created or updated").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("group", "Is the group of the kind")).
			Param(ws.PathParameter("version", "Is the version of the kind")).
			Param(ws.PathParameter("kind", "Is the kind of the resource")).
			Param(ws.PathParameter("namespace", "Is the namespace of the resource")).
			Param(ws.PathParameter("resource", "Is the name of the resource")).
			Param(ws.QueryParameter("source", "The producer of the alerting rule")).
			Param(ws.QueryParameter("status", "The alert to filter the results by")).
			Operation("ListResourceAlerts").
			Returns(http.StatusOK, "The alert has been successfully stored", monitoring.AlertList{}),
	)

	ws.Route(
		withAllErrors(ws.GET("/teams/{team}/alerts")).To(a.findAlertsByTeam).
			Filter(filters.DefaultMembersHandler.Filter).
			Doc("Returns all available alerts currently in kore").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("team", "Is the name of the team the alerts reside")).
			Param(ws.QueryParameter("history", "The number of historical records to retrieve")).
			Param(ws.QueryParameter("status", "The alert to filter the results by")).
			Param(ws.QueryParameter("latest", "Indicates to we only want the latest alert status")).
			Operation("ListTeamAlerts").
			Returns(http.StatusOK, "Listing of the alert statuses in kore", monitoring.AlertList{}),
	)

	ws.Route(
		withAllErrors(ws.GET("/teams/{team}/rules")).To(a.findRulesByTeam).
			Filter(filters.DefaultMembersHandler.Filter).
			Doc("Returns all available rules related to the team resources").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(ws.PathParameter("team", "Is the name of the team the alerts reside")).
			Param(ws.QueryParameter("history", "The number of historical records to retrieve")).
			Param(ws.QueryParameter("status", "The alert to filter the results by")).
			Operation("ListTeamRules").
			Returns(http.StatusOK, "A list of the rules", monitoring.AlertList{}),
	)

	return ws, nil
}

func (a *alertsHandler) findAlertsOnResource(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()

		options := kore.ListAlertOptions{
			Resource: &corev1.Ownership{
				Group:     req.PathParameter("group"),
				Version:   req.PathParameter("version"),
				Kind:      req.PathParameter("kind"),
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("resource"),
			},
			Source:   req.QueryParameter("source"),
			Statuses: req.QueryParameters("status"),
		}

		model, err := a.AlertRules().ListAlerts(ctx, options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

func (a *alertsHandler) updateAlert(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		model := &monitoring.Alert{}

		if err := req.ReadEntity(model); err != nil {
			return err
		}

		if err := a.AlertRules().UpdateAlert(ctx, model); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

func (a *alertsHandler) purgeRule(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")

		keep := req.QueryParameter("duration")
		if keep == "" {
			keep = "24h"
		}

		duration, err := time.ParseDuration(keep)
		if err != nil {
			return err
		}

		rule, err := a.AlertRules().GetRule(ctx, name,
			corev1.Ownership{
				Group:     req.PathParameter("group"),
				Version:   req.PathParameter("version"),
				Kind:      req.PathParameter("kind"),
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("resource"),
			},
		)
		if err != nil {
			return err
		}

		if err := a.AlertRules().PurgeHistory(ctx, rule, duration); err != nil {
			return err
		}

		return nil
	})
}

func (a *alertsHandler) findRule(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")

		rule, err := a.AlertRules().GetRule(ctx, name,
			corev1.Ownership{
				Group:     req.PathParameter("group"),
				Version:   req.PathParameter("version"),
				Kind:      req.PathParameter("kind"),
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("resource"),
			},
		)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, rule)
	})
}

func (a *alertsHandler) findResourceRules(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()

		options := kore.ListAlertOptions{
			Resource: &corev1.Ownership{
				Group:     req.PathParameter("group"),
				Version:   req.PathParameter("version"),
				Kind:      req.PathParameter("kind"),
				Namespace: req.PathParameter("namespace"),
				Name:      req.PathParameter("resource"),
			},
			Source:   req.QueryParameter("source"),
			Statuses: req.QueryParameters("status"),
		}

		if req.QueryParameter("latest") == "true" {
			options.History = 0
		}

		if req.QueryParameter("history") != "" {
			count, err := strconv.ParseInt(req.QueryParameter("history"), 10, 64)
			if err != nil {
				return err
			}
			options.History = int(count)
		}

		model, err := a.AlertRules().ListRules(ctx, options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

func (a *alertsHandler) deleteResourceRules(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()

		return a.AlertRules().DeleteResourceRules(ctx, corev1.Ownership{
			Group:     req.PathParameter("group"),
			Version:   req.PathParameter("version"),
			Kind:      req.PathParameter("kind"),
			Namespace: req.PathParameter("namespace"),
			Name:      req.PathParameter("resource"),
		}, req.QueryParameter("source"))
	})
}

func (a *alertsHandler) deleteRule(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")

		model := &monitoring.Rule{}
		if err := req.ReadEntity(model); err != nil {
			return err
		}
		model.Name = name

		if _, err := a.AlertRules().DeleteRule(ctx, model); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

func (a *alertsHandler) findRulesByTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()

		options := kore.ListAlertOptions{
			Statuses: req.QueryParameters("status"),
			Team:     req.PathParameter("team"),
		}
		if req.QueryParameter("history") != "" {
			count, err := strconv.ParseInt(req.QueryParameter("history"), 10, 64)
			if err != nil {
				return err
			}
			options.History = int(count)
		}

		model, err := a.AlertRules().ListRules(ctx, options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

func (a *alertsHandler) findAlertsByTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()

		options := kore.ListAlertOptions{
			Statuses: req.QueryParameters("status"),
			Team:     req.PathParameter("team"),
		}

		if req.QueryParameter("history") != "" {
			count, err := strconv.ParseInt(req.QueryParameter("history"), 10, 64)
			if err != nil {
				return err
			}
			options.History = int(count)
		}

		model, err := a.AlertRules().ListAlerts(ctx, options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

// findAllRules is responsible for rendering all the alerts
func (a *alertsHandler) findAllRules(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		options := kore.ListAlertOptions{
			Source: req.QueryParameter("source"),
		}

		model, err := a.AlertRules().ListRules(req.Request.Context(), options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

// findAllAlerts is responsible for rendering all the alerts
func (a *alertsHandler) findAllAlerts(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		options := kore.ListAlertOptions{
			Labels:   req.QueryParameters("label"),
			Source:   req.QueryParameter("source"),
			Statuses: req.QueryParameters("status"),
		}

		if req.QueryParameter("latest") == "true" {
			options.History = 0
		}

		if req.QueryParameter("history") != "" {
			count, err := strconv.ParseInt(req.QueryParameter("history"), 10, 64)
			if err != nil {
				return err
			}
			options.History = int(count)
		}

		model, err := a.AlertRules().ListAlerts(req.Request.Context(), options)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, model)
	})
}

// updateRule is responsible for updating or creating a rule
func (a *alertsHandler) updateRule(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		model := &monitoring.Rule{}
		if err := req.ReadEntity(model); err != nil {
			return err
		}
		model.Name = req.PathParameter("name")

		entity, err := a.AlertRules().UpdateRule(req.Request.Context(), model)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, entity)
	})
}

// Name returns the name of the handler
func (a alertsHandler) Name() string {
	return "alerts"
}

// EnableAudit defaults to audit everything.
func (a alertsHandler) EnableAudit() bool {
	return false
}
