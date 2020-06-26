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
	"strings"
	"time"

	"github.com/appvia/kore/pkg/persistence/model"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Alerts provides access to the alerts
type Alerts interface {
	// Get retrieves a alert by filter
	Get(context.Context, ...ListFunc) (*model.Alert, error)
	// Delete removes an alert
	Delete(context.Context, *model.Alert) error
	// DeleteBy removes an alert by filter
	DeleteBy(context.Context, ...ListFunc) error
	// List returns a filtered list of alerts
	List(context.Context, ...ListFunc) ([]*model.Alert, error)
	// Preload allows for the consumer to select the preloaded fields
	Preload(...string) Alerts
	// PurgeHistory purges the history of alerts for a rule
	PurgeHistory(context.Context, time.Duration, ...ListFunc) error
	// Transaction set the db connection
	Transaction(*gorm.DB) Alerts
	// Update updates or creates an alert
	Update(context.Context, *model.Alert) error
}

type alertsImpl struct {
	Interface
	// load is the preloaded fields
	load []string
	// conn is the db connection for this query
	conn *gorm.DB
	// transaction indicates a tx has been set
	transaction bool
}

// Delete removes an alert
func (i *alertsImpl) Delete(ctx context.Context, iv *model.Alert) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return i.conn.Delete(iv).Error
}

// DeleteBy removes alerts by filter
func (i *alertsImpl) DeleteBy(ctx context.Context, filters ...ListFunc) error {
	if len(filters) <= 0 {
		return errors.New("no filters defined for deletion of alerts")
	}

	list, err := i.List(ctx, filters...)
	if err != nil {
		return err
	}

	for _, x := range list {
		if err := i.conn.Model(&model.Alert{}).Delete(x).Error; err != nil {
			return err
		}
	}

	return nil
}

// Get returns a alert by filter - else a no record found error
func (i *alertsImpl) Get(ctx context.Context, opts ...ListFunc) (*model.Alert, error) {
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

// List returns a filtered list of alerts
func (i *alertsImpl) List(ctx context.Context, opts ...ListFunc) ([]*model.Alert, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	q := Preload(i.load, i.conn).
		Preload("Labels").
		Preload("Rule").
		Preload("Rule.Team").
		Select("i.*").
		Table("alerts i").
		Joins("JOIN alert_rules r ON r.id = i.rule_id").
		Joins("JOIN teams t ON t.id = r.team_id")

	var statuses []string

	if terms.HasRuleID() {
		q = q.Where("i.rule_id = ?", terms.GetRuleID())
	}
	if terms.HasAlertStatus() {
		statuses = terms.GetAlertStatus()
	}
	if terms.HasStatus() {
		statuses = []string{terms.GetStatus()}
	}
	if len(statuses) > 0 {
		q = q.Where("i.status IN (?)", statuses)
	}
	if terms.HasAlertSeverity() {
		q = q.Where("r.severity = ?", terms.GetAlertSeverity())
	}
	if terms.HasGroup() {
		q = q.Where("r.resource_group = ?", terms.GetGroup())
	}
	if terms.HasVersion() {
		q = q.Where("r.resource_version = ?", terms.GetVersion())
	}
	if terms.HasKind() {
		q = q.Where("r.resource_kind = ?", terms.GetKind())
	}
	if terms.HasNamespace() {
		q = q.Where("r.resource_namespace = ?", terms.GetNamespace())
	}
	if terms.HasResourceName() {
		q = q.Where("r.resource_name = ?", terms.GetResourceName())
	}
	if terms.HasName() {
		q = q.Where("r.name = ?", terms.GetName())
	}
	if terms.HasID() {
		q = q.Where("i.id = ?", terms.GetID())
	}
	if terms.HasTeam() {
		q = q.Where("t.name = ?", terms.GetTeam())
	}
	if terms.HasTeamID() {
		q = q.Where("t.id = ?", terms.GetTeamID())
	}
	if terms.HasAlertLatest() {
		q = q.Where("archived_at IS NULL")
	}
	if terms.HasStatus() {
		q = q.Where("i.status IN (?)", statuses)
	}
	if terms.HasAlertSource() {
		q = q.Where("r.source = ?", terms.GetAlertSource())
	}
	if terms.HasAlertUID() {
		q = q.Where("u.uid = ?", terms.GetAlertUID())
	}
	if terms.HasAlertHistory() {
		q = q.Limit(terms.GetAlertHistory())
	}
	if terms.HasAlertFingerprint() {
		q = q.Where("i.fingerprint = ?", terms.GetAlertFingerprint())
	}

	var list []*model.Alert

	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}

	if terms.HasAlertLabels() {
		keypairs := make(map[string]string)
		for _, x := range terms.GetAlertLabels() {
			items := strings.Split(x, "=")
			keypairs[items[0]] = items[1]
		}
		filtered := []*model.Alert{}

		for _, x := range list {
			if x.HasLabels(keypairs) {
				filtered = append(filtered, x)
			}
		}

		return filtered, nil
	}

	return list, nil
}

// Update updates or creates an alert
func (i *alertsImpl) Update(ctx context.Context, iv *model.Alert) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	if iv.Rule != nil {
		iv.RuleID = iv.Rule.ID
	}
	if iv.RuleID == 0 {
		return errors.New("alert has no rule id")
	}

	updateFn := func(tx *gorm.DB) error {

		// @case: if we only have the only alert - i.e. the initial one we don't
		// need to worry, we just invalidate that and inject the new one
		list, err := i.Alerts().Transaction(tx).List(ctx, Filter.WithRuleID(iv.RuleID))
		if err != nil {
			return err
		}
		switch len(list) {
		case 0:
			panic("we have a rule without any status, which should never happen")
		case 1:
			if iv.Status == model.AlertStatusOK {
				return nil
			}

			if err = tx.
				Exec("UPDATE alerts SET archived_at = NOW() WHERE rule_id = ?", iv.RuleID).
				Error; err != nil {
				return err
			}

			return tx.
				Set("gorm:association_autoupdate", false).
				Save(iv).
				Error
		}

		rule, err := i.AlertRules().Transaction(tx).Get(ctx, Filter.WithID(iv.RuleID))
		if err != nil {
			return err
		}

		// @case we have an alert for this rule, we need to check if it's for the same
		// instance
		list, err = i.Alerts().Transaction(tx).List(ctx,
			Filter.WithAlertFingerprint(iv.Fingerprint),
			Filter.WithRuleID(rule.ID),
			Filter.WithAlertLatest(),
		)
		if err != nil {
			log.Warn("we didn't find any alerts on the rule")

			return err
		}
		switch len(list) {
		// we have a new instance of the same rule
		case 0:
			if err = tx.
				Exec("UPDATE alerts SET archived_at = NOW() WHERE rule_id = ? AND status = 'OK'", rule.ID).
				Error; err != nil {
				return err
			}

			return tx.
				Set("gorm:association_autoupdate", false).
				Save(iv).
				Error
		case 1:
			if list[0].Status == iv.Status {
				return nil
			}

			if err = tx.
				Exec("UPDATE alerts SET archived_at = NOW() WHERE archived_at IS NULL AND fingerprint = ? AND rule_id = ?",
					list[0].Fingerprint, list[0].RuleID).
				Error; err != nil {
				return err
			}

			return tx.
				Set("gorm:association_autoupdate", false).
				Save(iv).
				Error

		default:
			log.WithFields(log.Fields{
				"fingerprint": iv.Fingerprint,
				"rule":        iv.Rule.Name,
				"team":        iv.Rule.ResourceNamespace,
				"resource":    iv.Rule.ResourceName,
			}).Error("alert has multiple statuses")
		}

		return nil
	}

	if i.transaction {
		return updateFn(i.conn)
	}

	return i.conn.Transaction(updateFn)
}

// PurgeHistory purges the history of alerts for a rule
func (i *alertsImpl) PurgeHistory(ctx context.Context, keep time.Duration, filters ...ListFunc) error {
	return i.conn.Transaction(func(tx *gorm.DB) error {

		current, err := i.AlertRules().Get(ctx, filters...)
		if err != nil {
			return err
		}
		before := time.Now().Add(-keep).Format("2006-01-02 15:04:05")

		return tx.
			Exec("DELETE FROM alerts WHERE rule_id = ? AND created_at > ? AND archived_at IS NOT NULL", current.ID, before).
			Error
	})
}

// Preload allows for the consumer to select the preloaded fields
func (i *alertsImpl) Preload(v ...string) Alerts {
	i.load = append(i.load, v...)

	return i
}

// Transaction sets the connection transaction
func (i *alertsImpl) Transaction(tx *gorm.DB) Alerts {
	i.conn = tx
	i.transaction = true

	return i
}
