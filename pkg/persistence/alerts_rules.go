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

package persistence

import (
	"context"
	"errors"

	"github.com/appvia/kore/pkg/persistence/model"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
)

// AlertRules provides access to the alert rules
type AlertRules interface {
	// Get retrieves a alert rule by filter
	Get(context.Context, ...ListFunc) (*model.AlertRule, error)
	// Delete removes an alert rule
	Delete(context.Context, *model.AlertRule) error
	// DeleteBy removes an alert rule by filter
	DeleteBy(context.Context, ...ListFunc) error
	// List returns a filtered list of alert rule rules
	List(context.Context, ...ListFunc) ([]*model.AlertRule, error)
	// Preload allows for the consumer to select the preloaded fields
	Preload(...string) AlertRules
	// Transaction set the db transation
	Transaction(*gorm.DB) AlertRules
	// Update updates or creates and alert rule rules
	Update(context.Context, *model.AlertRule) error
	// UpdateStatus is responsible for update the rule status
	UpdateStatus(context.Context, *model.AlertRule, ...ListFunc) error
}

type arulesImpl struct {
	Interface
	// load is the preloaded fields
	load []string
	// conn is the db connection for this query
	conn *gorm.DB
	// istransaction
	istransation bool
}

// Delete removes an alert rule
func (i *arulesImpl) Delete(ctx context.Context, iv *model.AlertRule) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return i.conn.Delete(iv).Error
}

// DeleteBy removes alert rule rules by filter
func (i *arulesImpl) DeleteBy(ctx context.Context, filters ...ListFunc) error {
	if len(filters) <= 0 {
		return errors.New("no filters defined for deletion of alert rule rules")
	}

	list, err := i.makeListQuery(ctx, filters...)
	if err != nil {
		return err.Error
	}

	for x := 0; x < len(list); x++ {
		if err := i.conn.Model(&model.AlertRule{}).Delete(&list[x]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Get returns a alert rule by filter - else a no record found error
func (i *arulesImpl) Get(ctx context.Context, opts ...ListFunc) (*model.AlertRule, error) {
	list, err := i.List(ctx, opts...)
	if err != nil {
		return nil, err
	}
	switch len(list) {
	case 0:
		return nil, gorm.ErrRecordNotFound
	case 1:
		return list[0], nil
	default:
		return nil, errors.New("matched more than one record")
	}
}

// List returns a filtered list of alert rule rules
func (i *arulesImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.AlertRule, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	list, q := i.makeListQuery(ctx, opts...)

	return list, q.Error
}

// UpdateStatus updates the status or alert history for a given rule
func (i *arulesImpl) UpdateStatus(ctx context.Context, rule *model.AlertRule, filters ...ListFunc) error {
	switch len(rule.Alerts) {
	case 1:
	case 0:
		return errors.New("rule contains no alert status")
	default:
		return errors.New("rule update contains multiple statuses")
	}

	return i.conn.Transaction(func(tx *gorm.DB) error {
		current, err := i.AlertRules().Transaction(tx).Get(ctx, filters...)
		if err != nil {
			return err
		}
		alert := &rule.Alerts[0]
		alert.Rule = current

		return i.Alerts().Transaction(tx).Update(ctx, alert)
	})
}

// Update updates or creates and alert rule
func (i *arulesImpl) Update(ctx context.Context, iv *model.AlertRule) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	if iv.Team != nil {
		iv.TeamID = iv.Team.ID
	}
	if iv.TeamID == 0 {
		return errors.New("no team id defined")
	}

	updateFn := func(tx *gorm.DB) error {
		entity := &model.AlertRule{}

		err := tx.
			Where("name = ?", iv.Name).
			Where("resource_group = ? AND resource_version = ? AND resource_kind = ?", iv.ResourceGroup, iv.ResourceVersion, iv.ResourceKind).
			Where("resource_namespace = ? and resource_name = ?", iv.ResourceNamespace, iv.ResourceName).
			Where("source = ?", iv.Source).
			First(entity).
			Error

		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				err := tx.
					Set("gorm:save_associations", false).
					Save(iv).
					Error
				if err != nil {
					return err
				}

				return tx.
					Set("gorm:save_associations", false).
					Save(&model.Alert{
						RuleID:  iv.ID,
						Summary: "none",
						Status:  model.AlertStatusOK,
					}).
					Error
			}

			return err
		}
		iv.ID = entity.ID

		return tx.
			Where("id = ?", iv.ID).
			Save(iv).
			Error
	}

	if i.istransation {
		return updateFn(i.conn)
	}

	return i.conn.Transaction(updateFn)
}

// Preload allows for the consumer to select the preloaded fields
func (i *arulesImpl) Preload(v ...string) AlertRules {
	i.load = append(i.load, v...)

	return i
}

func (i *arulesImpl) makeListQuery(ctx context.Context, opts ...ListFunc) ([]*model.AlertRule, *gorm.DB) {
	terms := ApplyListOptions(opts...)

	list := []*model.AlertRule{}

	q := Preload(i.load, i.conn).
		Preload("Team").
		Select("i.*").
		Table("alert_rules i").
		Joins("LEFT JOIN teams t ON t.id = i.team_id").
		Joins("LEFT JOIN alerts a ON a.rule_id = i.id").
		Where("a.archived_at IS NULL").
		Group("i.id").
		Order("a.created_at DESC")

	switch {
	case terms.HasAlertHistory():
		q = q.Preload("Alerts", func(db *gorm.DB) *gorm.DB {
			return db.Limit(terms.GetAlertHistory()).Order("created_at DESC")
		}).Preload("Alerts.Labels").Find(&list)

	case terms.HasAlertLatest():
		q = q.Preload("Alerts", func(db *gorm.DB) *gorm.DB {
			return db.Where("archived_at IS NULL").Order("created_at DESC")
		}).Preload("Alerts.Labels").Find(&list)

	case terms.HasStatus():
		q = q.Preload("Alerts", func(db *gorm.DB) *gorm.DB {
			return db.Where("status IN (?)", terms.GetStatus()).Order("created_at DESC")
		}).Preload("Alerts.Labels").Find(&list)

	default:
		q = q.Preload("Alerts", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).Preload("Alerts.Labels").Find(&list)
	}
	if terms.HasStatus() {
		q = q.Where("a.status = ?", terms.GetStatus())
	}
	if terms.HasAlertStatus() {
		q = q.Where("a.status IN (?)", terms.GetAlertStatus())
	}
	if terms.HasAlertSource() {
		q = q.Where("i.source = ?", terms.GetAlertSource())
	}
	if terms.HasAlertSeverity() {
		q = q.Where("i.severity = ?", terms.GetAlertSeverity())
	}
	if terms.HasName() {
		q = q.Where("i.name = ?", terms.GetName())
	}
	if terms.HasTeam() {
		q = q.Where("t.name = ?", terms.GetTeam())
	}
	if terms.HasTeamID() {
		q = q.Where("t.id = ?", terms.GetTeamID())
	}
	if terms.HasGroup() {
		q = q.Where("resource_group = ?", terms.GetGroup())
	}
	if terms.HasVersion() {
		q = q.Where("resource_version = ?", terms.GetVersion())
	}
	if terms.HasKind() {
		q = q.Where("resource_kind = ?", terms.GetKind())
	}
	if terms.HasNamespace() {
		q = q.Where("resource_namespace = ?", terms.GetNamespace())
	}
	if terms.HasResourceName() {
		q = q.Where("resource_name = ?", terms.GetResourceName())
	}
	if terms.HasID() {
		q = q.Where("i.id = ?", terms.GetID())
	}

	return list, q.Find(&list)
}

func (i *arulesImpl) Transaction(tx *gorm.DB) AlertRules {
	i.conn = tx
	i.istransation = true

	return i
}
