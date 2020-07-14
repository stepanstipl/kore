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

package kore

import (
	"context"
	"regexp"
	"strings"
	"time"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	monitoring "github.com/appvia/kore/pkg/apis/monitoring/v1beta1"
	"github.com/appvia/kore/pkg/kore/authentication"
	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/validation"

	log "github.com/sirupsen/logrus"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
)

var _ AlertRules = &alertsImpl{}

var (
	// AlertStatuses is a collection of possible alert statuses
	AlertStatuses = []string{
		monitoring.AlertStatusActive,
		monitoring.AlertStatusOK,
		monitoring.AlertStatusSilenced,
	}

	// AlertSeverity is a collection severity we support
	AlertSeverity = []string{
		"critical",
		"warning",
		"ok",
		"none",
	}
)

// ListAlertOptions are options to the listing
type ListAlertOptions struct {
	// Name is the name of the rule
	Name string
	// History is the number of historical records to retrieve.
	// Default to latest, you can use -1 to return all
	History int
	// Labels is a collection of labels to filter by
	Labels []string
	// Ownership references the alerts a specific resource
	Resource *corev1.Ownership
	// Severity is the severity we are looking for
	Severity string
	// Statuses is the status we are looking for
	Statuses []string
	// Source is the provider of the rule
	Source string
	// Team is the team to we should look in
	Team string
}

// AlertRules provides the interface for the alert rules
type AlertRules interface {
	// ListAlerts returns all the alerts in kore
	ListAlerts(ctx context.Context, options ListAlertOptions) (*monitoring.AlertList, error)
	// ListRules returns all the rules in kore
	ListRules(ctx context.Context, options ListAlertOptions) (*monitoring.AlertRuleList, error)
	// DeleteRule is responsible for deleting a rule on a resource and all alert statuses thereafter
	DeleteRule(ctx context.Context, rule *monitoring.AlertRule) (*monitoring.AlertRule, error)
	// DeleteResourceRules is responsible for deleting all rules on a resource
	DeleteResourceRules(ctx context.Context, resource corev1.Ownership, source string) error
	// GetRule returns a specific rule
	GetRule(ctx context.Context, name string, resource corev1.Ownership) (*monitoring.AlertRule, error)
	// PurgeHistory removes the history an a rule
	PurgeHistory(ctx context.Context, rule *monitoring.AlertRule, keep time.Duration) error
	// SilenceRule is responsible for suspending any alerts on a rule
	SilenceRule(ctx context.Context, rule monitoring.AlertRule, message string, duration time.Duration) error
	// UpdateRule is responsible for updating or creating a rule on a resource
	UpdateRule(ctx context.Context, rule *monitoring.AlertRule) (*monitoring.AlertRule, error)
	// UpdateAlert is responsible for updating the status of a rule
	UpdateAlert(ctx context.Context, rule *monitoring.Alert) error
}

type alertsImpl struct {
	Interface
}

var (
	alertLabelRegex = regexp.MustCompile(`^[a-zA-Z\-/].*=[a-zA-Z\-/].*$`)
	sourceNameRegex = regexp.MustCompile(`^([a-zA-Z0-9]){3,32}$`)
)

// ListAlerts returns all the alerts in kore
func (a *alertsImpl) ListAlerts(ctx context.Context, options ListAlertOptions) (*monitoring.AlertList, error) {
	user := authentication.MustGetIdentity(ctx)

	// @step: input validation
	if !user.IsGlobalAdmin() {
		if options.Team == "" {
			return nil, NewErrNotAllowed("You don't have permissions to these resources")
		}
		if options.Team != "" && !user.IsMember(options.Team) {
			return nil, NewErrNotAllowed("You must be a member of the team")
		}
		if options.Resource != nil && !user.IsMember(options.Team) {
			return nil, NewErrNotAllowed("You must be a member of the team")
		}
	}

	for _, x := range options.Statuses {
		if !utils.Contains(x, model.AlertStatuses) {
			return nil, validation.NewError("invalid status").WithFieldErrorf(
				"statuses", validation.InvalidValue, "status: %s does not exist", x)
		}
	}

	for _, x := range options.Labels {
		if !alertLabelRegex.MatchString(x) {
			return nil, validation.NewError("invalid label").WithFieldErrorf(
				"labels", validation.InvalidValue, "label: %s is invalid", x)
		}
	}

	var filters []persistence.ListFunc

	// @step: apply the filters
	if options.Team != "" {
		filters = append(filters, persistence.Filter.WithTeam(options.Team))
	}
	if options.Name != "" {
		filters = append(filters, persistence.Filter.WithName(options.Name))
	}
	if options.Resource != nil {
		filters = append(filters, []persistence.ListFunc{
			persistence.Filter.WithNamespace(options.Resource.Namespace),
			persistence.Filter.WithResourceGroup(options.Resource.Group),
			persistence.Filter.WithResourceKind(options.Resource.Kind),
			persistence.Filter.WithResourceVersion(options.Resource.Version),
			persistence.Filter.WithResourceName(options.Resource.Name),
		}...)
	}
	if len(options.Statuses) > 0 {
		filters = append(filters, persistence.Filter.WithAlertStatus(options.Statuses))
	}
	if options.Source != "" {
		filters = append(filters, persistence.Filter.WithAlertSource(options.Source))
	}
	if options.History == 0 {
		filters = append(filters, persistence.Filter.WithAlertLatest())
	} else {
		filters = append(filters, persistence.Filter.WithAlertHistory(options.History))
	}
	if len(options.Labels) > 0 {
		filters = append(filters, persistence.Filter.WithAlertLabels(options.Labels))
	}

	list, err := a.Persist().Alerts().List(ctx, filters...)
	if err != nil {
		return nil, err
	}

	return DefaultConvertor.FromAlertsModelList(list)
}

// GetRule returns a specific rule
func (a *alertsImpl) GetRule(ctx context.Context, name string, resource corev1.Ownership) (*monitoring.AlertRule, error) {
	// @step: check the ownership
	if err := IsOwnershipValid(resource); err != nil {
		return nil, validation.NewError("rule has failed validation").WithFieldError(
			"resource", validation.InvalidValue, err.Error())
	}

	if err := a.userPermitted(ctx, resource); err != nil {
		return nil, ErrNotAllowed{message: "not permitted to update alert on this resource"}
	}

	entity, err := a.Persist().AlertRules().Get(ctx,
		persistence.Filter.WithName(name),
		persistence.Filter.WithNamespace(resource.Namespace),
		persistence.Filter.WithResourceGroup(resource.Group),
		persistence.Filter.WithResourceKind(resource.Kind),
		persistence.Filter.WithResourceVersion(resource.Version),
		persistence.Filter.WithResourceName(resource.Name),
	)
	if err != nil {
		if persistence.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return DefaultConvertor.FromAlertRuleModel(entity), nil
}

// ListRules returns the rules in kore
func (a *alertsImpl) ListRules(ctx context.Context, options ListAlertOptions) (*monitoring.AlertRuleList, error) {
	user := authentication.MustGetIdentity(ctx)

	// @step: input validation
	if !user.IsGlobalAdmin() {
		if options.Team == "" {
			return nil, NewErrNotAllowed("You don't have permissions to these resources")
		}
		if options.Team != "" && !user.IsMember(options.Team) {
			return nil, NewErrNotAllowed("You must be a member of the team")
		}
		if options.Resource != nil && !user.IsMember(options.Team) {
			return nil, NewErrNotAllowed("You must be a member of the team")
		}
	}

	var filters []persistence.ListFunc

	// @step: apply the filters
	if options.Name != "" {
		filters = append(filters, persistence.Filter.WithName(options.Name))
	}
	if options.Team != "" {
		filters = append(filters, persistence.Filter.WithTeam(options.Team))
	}
	if options.Severity != "" {
		filters = append(filters, persistence.Filter.WithAlertSeverity(options.Severity))
	}
	if options.Resource != nil {
		filters = append(filters, []persistence.ListFunc{
			persistence.Filter.WithNamespace(options.Resource.Namespace),
			persistence.Filter.WithResourceGroup(options.Resource.Group),
			persistence.Filter.WithResourceKind(options.Resource.Kind),
			persistence.Filter.WithResourceVersion(options.Resource.Version),
			persistence.Filter.WithResourceName(options.Resource.Name),
		}...)
	}
	if options.Source != "" {
		filters = append(filters, persistence.Filter.WithAlertSource(options.Source))
	}

	list, err := a.Persist().AlertRules().List(ctx, filters...)
	if err != nil {
		return nil, err
	}

	return DefaultConvertor.FromAlertsRuleModelList(list)
}

// DeleteRule is responsible for deleting a rule on a resource and all alert statuses thereafter
func (a *alertsImpl) DeleteRule(ctx context.Context, rule *monitoring.AlertRule) (*monitoring.AlertRule, error) {
	// @step: ensure the ownership is valid
	if err := IsOwnershipValid(rule.Spec.Resource); err != nil {
		return nil, validation.NewError("rule has failed validation").WithFieldError(
			"spec.resource", validation.MustExist, err.Error())
	}

	// @step: ensure the user has access to the resources
	if err := a.userPermitted(ctx, rule.Spec.Resource); err != nil {
		return nil, err
	}

	filters := []persistence.ListFunc{
		persistence.Filter.WithName(rule.Name),
		persistence.Filter.WithNamespace(rule.Spec.Resource.Namespace),
		persistence.Filter.WithResourceGroup(rule.Spec.Resource.Group),
		persistence.Filter.WithResourceKind(rule.Spec.Resource.Kind),
		persistence.Filter.WithResourceVersion(rule.Spec.Resource.Version),
		persistence.Filter.WithResourceName(rule.Spec.Resource.Name),
	}

	entity, err := a.Persist().AlertRules().Get(ctx, filters...)
	if err != nil {
		if persistence.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return DefaultConvertor.FromAlertRuleModel(entity), a.Persist().AlertRules().Delete(ctx, entity)
}

// DeleteResourceRules is responsible for deleting all rules on a resource
func (a *alertsImpl) DeleteResourceRules(ctx context.Context, resource corev1.Ownership, source string) error {
	logger := log.WithFields(log.Fields{
		"group":   resource.Group,
		"kind":    resource.Kind,
		"name":    resource.Name,
		"team":    resource.Namespace,
		"version": resource.Version,
	})
	logger.Info("attempting to delete all the rules on the resource")

	// @step: input validation
	if err := IsOwnershipValid(resource); err != nil {
		return ErrNotAllowed{message: "invalid resource ownership"}
	}
	if !schema.GetScheme().Recognizes(kschema.GroupVersionKind{
		Group:   resource.Group,
		Version: resource.Version,
		Kind:    resource.Kind,
	}) {
		return ErrNotAllowed{message: "resource type not found"}
	}
	if !ValidateSourceName(source) {
		return ErrNotAllowed{message: "invalid source name"}
	}
	if err := a.userPermitted(ctx, resource); err != nil {
		return err
	}

	return a.Persist().AlertRules().DeleteBy(ctx,
		persistence.Filter.WithAlertSource(source),
		persistence.Filter.WithNamespace(resource.Namespace),
		persistence.Filter.WithResourceGroup(resource.Group),
		persistence.Filter.WithResourceKind(resource.Kind),
		persistence.Filter.WithResourceName(resource.Name),
		persistence.Filter.WithResourceVersion(resource.Version),
		persistence.Filter.WithTeam(resource.Namespace),
	)
}

// SilenceRule is responsible for suspending any alerts on a rule
func (a *alertsImpl) SilenceRule(ctx context.Context, rule monitoring.AlertRule, message string, duration time.Duration) error {
	if message == "" {
		return ErrNotAllowed{message: "you must specify a reason for silence"}
	}
	if duration < 0 {
		return ErrNotAllowed{message: "duration must be positive"}
	}

	if err := IsOwnershipValid(rule.Spec.Resource); err != nil {
		return ErrNotAllowed{message: "invalid resource ownership"}
	}

	// @step: ensure the user has access to the resources
	if err := a.userPermitted(ctx, rule.Spec.Resource); err != nil {
		return err
	}

	logger := log.WithFields(log.Fields{
		"kind": rule.Spec.Resource.Kind,
		"name": rule.Spec.Resource.Name,
		"team": rule.Spec.Resource.Namespace,
	})
	logger.Info("attempting to silence the alert")

	entity, err := a.Persist().AlertRules().Get(ctx,
		persistence.Filter.WithResourceGroup(rule.Spec.Resource.Group),
		persistence.Filter.WithResourceVersion(rule.Spec.Resource.Version),
		persistence.Filter.WithResourceKind(rule.Spec.Resource.Kind),
		persistence.Filter.WithNamespace(rule.Spec.Resource.Namespace),
		persistence.Filter.WithResourceName(rule.Spec.Resource.Name),
		persistence.Filter.WithAlertLatest(),
	)
	if err != nil {
		if persistence.IsNotFound(err) {
			return ErrNotFound
		}

		return err
	}

	if entity.Alerts[0].Status == model.AlertStatusSilenced {
		return nil
	}

	// @step: we create an new status for the rule
	status := &model.Alert{Rule: entity, Status: model.AlertStatusSilenced}

	return a.Persist().Alerts().Update(ctx, status)
}

// UpdateRule is responsible for updating or creating a rule on a resource
func (a *alertsImpl) UpdateRule(ctx context.Context, rule *monitoring.AlertRule) (*monitoring.AlertRule, error) {
	logger := log.WithFields(log.Fields{
		"group":    rule.Spec.Resource.Group,
		"name":     rule.Name,
		"resource": rule.Spec.Resource.Kind,
	})
	logger.Debug("attempting to update the monitoring rule")

	if err := a.ValidateRule(ctx, rule); err != nil {
		logger.WithError(err).Error("rule has failed validation")

		return nil, err
	}

	team, err := a.Persist().Teams().Get(ctx, rule.Spec.Resource.Namespace)
	if err != nil {
		return nil, err
	}

	m := DefaultConvertor.ToAlertRule(rule)
	m.Team = team

	filters := []persistence.ListFunc{
		persistence.Filter.WithName(rule.Name),
		persistence.Filter.WithNamespace(rule.Spec.Resource.Namespace),
		persistence.Filter.WithResourceGroup(rule.Spec.Resource.Group),
		persistence.Filter.WithResourceKind(rule.Spec.Resource.Kind),
		persistence.Filter.WithResourceVersion(rule.Spec.Resource.Version),
		persistence.Filter.WithResourceName(rule.Spec.Resource.Name),
	}

	if err := a.Persist().AlertRules().Update(ctx, m); err != nil {
		logger.WithError(err).Error("trying to update the rule")

		return nil, err
	}

	entity, err := a.Persist().AlertRules().Get(ctx, filters...)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the updated rule")

		return nil, err
	}

	return DefaultConvertor.FromAlertRuleModel(entity), nil
}

// PurgeHistory is responsible for purging the history of an alert
func (a *alertsImpl) PurgeHistory(ctx context.Context, rule *monitoring.AlertRule, keep time.Duration) error {
	logger := log.WithFields(log.Fields{
		"name":      rule.Name,
		"namespace": rule.Namespace,
		"resource":  rule.Spec.Resource.Kind,
	})
	logger.Debug("attempting to purge alert history for rules")

	// @step: check the ownership
	if err := IsOwnershipValid(rule.Spec.Resource); err != nil {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.resource", validation.InvalidValue, err.Error())
	}

	if err := a.userPermitted(ctx, rule.Spec.Resource); err != nil {
		return ErrNotAllowed{message: "not permitted to update alert on this resource"}
	}

	filter := []persistence.ListFunc{
		persistence.Filter.WithName(rule.Name),
		persistence.Filter.WithNamespace(rule.Spec.Resource.Namespace),
		persistence.Filter.WithResourceGroup(rule.Spec.Resource.Group),
		persistence.Filter.WithResourceKind(rule.Spec.Resource.Kind),
		persistence.Filter.WithResourceVersion(rule.Spec.Resource.Version),
		persistence.Filter.WithResourceName(rule.Spec.Resource.Name),
	}

	if err := a.Persist().Alerts().PurgeHistory(ctx, keep, filter...); err != nil {
		if persistence.IsNotFound(err) {
			return ErrNotFound
		}
		logger.WithError(err).Error("trying to purge history")

		return err
	}

	return nil
}

// UpdateAlert is responsible for updating the status of a rule
func (a *alertsImpl) UpdateAlert(ctx context.Context, alert *monitoring.Alert) error {
	logger := log.WithFields(log.Fields{
		"name":      alert.Name,
		"namespace": alert.Namespace,
	})
	logger.Debug("attempting to update the status of the rule")

	// @step: check the inputs
	if err := a.ValidateAlert(ctx, alert); err != nil {
		logger.WithError(err).Error("alert update has failed validation")

		return err
	}

	// @step: ensure the rule exists
	rule := alert.Status.Rule
	filters := []persistence.ListFunc{
		persistence.Filter.WithAlertSource(rule.Spec.Source),
		persistence.Filter.WithName(rule.Name),
		persistence.Filter.WithNamespace(rule.Spec.Resource.Namespace),
		persistence.Filter.WithResourceGroup(rule.Spec.Resource.Group),
		persistence.Filter.WithResourceKind(rule.Spec.Resource.Kind),
		persistence.Filter.WithResourceName(rule.Spec.Resource.Name),
		persistence.Filter.WithResourceVersion(rule.Spec.Resource.Version),
	}

	rulemodel, err := a.Persist().AlertRules().Get(ctx, filters...)
	if err != nil {
		logger.WithError(err).Error("trying to retrieve the alerting rule")

		if persistence.IsNotFound(err) {
			return ErrNotFound
		}

		return err
	}
	model := DefaultConvertor.ToAlert(alert)
	model.Rule = rulemodel

	// @step: we always make sure there's a rule to update first of all
	// @step: update the status of that rule is required
	if err := a.Persist().Alerts().Update(ctx, model); err != nil {
		if persistence.IsNotFound(err) {
			return ErrNotFound
		}
		logger.WithError(err).Error("trying to update the rule status")

		return err
	}

	return nil
}

func (a *alertsImpl) userPermitted(ctx context.Context, resource corev1.Ownership) error {
	user := authentication.MustGetIdentity(ctx)
	if user.IsGlobalAdmin() {
		return nil
	}

	resourceNamespace := resource.Namespace
	if resourceNamespace == "" {
		log.WithFields(log.Fields{
			"name":     resource,
			"resource": resource.Kind,
			"username": user.Username(),
		}).Warn("user trying to access the alerts not permitted")

		return NewErrNotAllowed("rule has not team associated")
	}

	if !user.IsMember(resourceNamespace) {
		log.WithFields(log.Fields{
			"name":     resource,
			"resource": resource.Kind,
			"username": user.Username(),
		}).Warn("user trying to access alerts not permitted")

		return NewErrNotAllowed("must be global admin or in team which owns this resource")
	}

	return nil
}

// ValidateSourceName checks the source name
func ValidateSourceName(name string) bool {
	return sourceNameRegex.MatchString(name)
}

// ValidateRule is used to check the inputs on a rule
func (a *alertsImpl) ValidateRule(ctx context.Context, o *monitoring.AlertRule) error {
	if o.Name == "" {
		return validation.NewError("rule failed validation").WithFieldError(
			"metadata.name", validation.Required, "name must exists")
	}
	if o.Namespace == "" {
		return validation.NewError("rule failed validation").WithFieldError(
			"metadata.namespace", validation.Required, "namespace must be defined")
	}
	if o.Spec.RawRule == "" {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.rawRule", validation.MustExist, "rule must have a raw definition")
	}
	if o.Spec.Summary == "" {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.summary", validation.MustExist, "rule must have a summary")
	}
	if o.Spec.Severity == "" {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.severity", validation.MustExist, "rule must have a severity")
	}
	if !utils.Contains(strings.ToLower(o.Spec.Severity), AlertSeverity) {
		return validation.NewError("rule has failed validation").WithFieldErrorf(
			"spec.severity", validation.MustExist, "rule severity %q is invalid", o.Spec.Severity)
	}
	if o.Spec.Source == "" {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.source", validation.MustExist, "rule must have a source")
	}
	if !ValidateSourceName(o.Spec.Source) {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.source", validation.InvalidValue, "rule source invalid")
	}
	if err := IsOwnershipValid(o.Spec.Resource); err != nil {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.resource", validation.InvalidValue, err.Error())
	}
	if !schema.GetScheme().Recognizes(kschema.GroupVersionKind{
		Group:   o.Spec.Resource.Group,
		Version: o.Spec.Resource.Version,
		Kind:    o.Spec.Resource.Kind,
	}) {
		return validation.NewError("rule has failed validation").WithFieldError(
			"spec.resource", validation.InvalidValue, "unknown resource type")
	}
	_, err := a.Persist().Teams().Get(ctx, o.Spec.Resource.Namespace)
	if err != nil {
		return err
	}
	if err := a.userPermitted(ctx, o.Spec.Resource); err != nil {
		return ErrNotAllowed{message: "not permitted to update alert on this resource"}
	}

	return nil
}

// ValidateAlert is used to check the inputs of an alert
func (a *alertsImpl) ValidateAlert(ctx context.Context, o *monitoring.Alert) error {
	if o.Name == "" {
		return validation.NewError("alert failed validation").WithFieldError(
			"metadata.name", validation.Required, "name must exists")
	}
	if o.Namespace == "" {
		return validation.NewError("alert failed validation").WithFieldError(
			"metadata.namespace", validation.Required, "namespace must be defined")
	}
	if o.GetAnnotations()["fingerprint"] == "" {
		return validation.NewError("alert failed validation").WithFieldError(
			"metadata.annotations['fingerprint']", validation.Required, "fingerprint must be defined")
	}
	if o.Spec.Event == "" {
		return validation.NewError("alert failed validation").WithFieldError(
			"spec.event", validation.Required, "alert must have an event")
	}
	if o.Status.Status == "" {
		return validation.NewError("alert failed validation").WithFieldError(
			"status.status", validation.Required, "alert must have a status")
	}
	if !utils.Contains(o.Status.Status, AlertStatuses) {
		return validation.NewError("alert failed validation").WithFieldError(
			"status.status", validation.InvalidValue, "invalid alert status")
	}
	if o.Status.Rule == nil {
		return validation.NewError("alert failed validation").WithFieldError(
			"status.rule", validation.Required, "alert must have a rule")
	}
	if err := a.ValidateRule(ctx, o.Status.Rule); err != nil {
		return validation.NewError("alert failed validation").WithFieldError(
			"status.rule", validation.InvalidValue, err.Error())
	}

	return nil
}
