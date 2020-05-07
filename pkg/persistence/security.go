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
	"time"

	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/appvia/kore/pkg/persistence/model"
)

// Security defines the security interface to the store
type Security interface {
	// GetScan returns a scan result from the store, with the rule results populated
	GetScan(context.Context, uint64) (*model.SecurityScanResult, error)
	// ListScans returns a list of scan results from the store, without the rule results populated
	ListScans(ctx context.Context, latestOnly bool, opts ...ListFunc) ([]*model.SecurityScanResult, error)
	// GetLatestResourceScan returns the latest scan for a specific resource
	GetLatestResourceScan(ctx context.Context, group string, version string, kind string, namespace string, name string) (*model.SecurityScanResult, error)
	// ListResourceScanHistory returns all scans for a specific resource
	ListResourceScanHistory(ctx context.Context, group string, version string, kind string, namespace string, name string) ([]*model.SecurityScanResult, error)
	// StoreScan stores the result of a scan in the store. If the supplied result has a zero archived_at, any
	// previous scan with a zero archived_at for the same resource will be compared to this scan (results, messages, overall status)
	// and if different, the old result will have its archived_at set, else the old result will be updated with an updated checked_at
	// time. If archived_at is set, it will simply be persisted with no changes to any other results.
	StoreScan(context.Context, *model.SecurityScanResult) error
	// ArchiveResourceScans sets any non-archived scan results for the resource to archived. This is for when a resource is deleted.
	ArchiveResourceScans(ctx context.Context, group string, version string, kind string, namespace string, name string) error
	// GetOverview retrieves an overview of the current security situation of all resources
	GetOverview(context.Context) (*model.SecurityOverview, error)
	// GetOverview retrieves an overview of the current security situation of all resources owned by a specific team
	GetTeamOverview(context.Context, string) (*model.SecurityOverview, error)
}

var _ Security = &securityImpl{}

type securityImpl struct {
	Interface
	// conn is the db connection
	conn *gorm.DB
}

func (s *securityImpl) GetScan(ctx context.Context, scanID uint64) (*model.SecurityScanResult, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	scan := &model.SecurityScanResult{}
	err := s.conn.Preload("Results").First(&scan, scanID).Error
	if err != nil {
		return nil, err
	}
	return scan, nil
}

func (s *securityImpl) ListScans(ctx context.Context, latestOnly bool, opts ...ListFunc) ([]*model.SecurityScanResult, error) {
	timed := prometheus.NewTimer(listLatency)
	defer timed.ObserveDuration()

	terms := ApplyListOptions(opts...)

	q := s.conn
	if latestOnly {
		q = q.Where("archived_at is null")
	}

	// If we have a full identity, use that, else use the individual terms
	if terms.HasIdentity() {
		q = q.Where(
			"resource_group = ? AND resource_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
			terms.GetGroup(),
			terms.GetVersion(),
			terms.GetKind(),
			terms.GetNamespace(),
			terms.GetName(),
		)
	} else {
		if terms.HasTeam() {
			q = q.Where("owning_team = ?", terms.GetTeam())
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
		if terms.HasName() {
			q = q.Where("resource_name = ?", terms.GetName())
		}
	}

	var list []*model.SecurityScanResult
	err := q.Find(&list).Error
	return list, err
}

func (s *securityImpl) GetLatestResourceScan(ctx context.Context, group string, version string, kind string, namespace string, name string) (*model.SecurityScanResult, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	scan, err := s.getLatestResourceScan(s.conn, group, version, kind, namespace, name)
	if err != nil {
		return nil, err
	}
	return scan, nil
}

func (s *securityImpl) getLatestResourceScan(tx *gorm.DB, group string, version string, kind string, namespace string, name string) (*model.SecurityScanResult, error) {
	scan := &model.SecurityScanResult{}
	err := s.conn.Preload("Results").
		Where("resource_group = ?", group).
		Where("resource_version = ?", version).
		Where("resource_kind = ?", kind).
		Where("resource_namespace = ?", namespace).
		Where("resource_name = ?", name).
		Where("archived_at IS NULL").
		First(&scan).
		Error
	if err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return scan, nil
}

func (s *securityImpl) ListResourceScanHistory(ctx context.Context, group string, version string, kind string, namespace string, name string) ([]*model.SecurityScanResult, error) {
	// This is just a convenience wrapper for the generic list scans function:
	return s.ListScans(ctx, false, Filter.WithIdentity(group, version, kind, namespace, name))
}

func (s *securityImpl) ArchiveResourceScans(ctx context.Context, group string, version string, kind string, namespace string, name string) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	return s.conn.Exec(
		"UPDATE security_scan_results SET archived_at = ? WHERE archived_at IS NULL AND resource_group = ? AND resource_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
		time.Now(),
		group,
		version,
		kind,
		namespace,
		name,
	).Error
}

func (s *securityImpl) StoreScan(ctx context.Context, result *model.SecurityScanResult) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	// Ensure the checking, archiving and inserting is done atomically in a single transaction.
	return s.conn.Transaction(func(tx *gorm.DB) error {
		// If archived at is set, just persist this without updating/checking current result, as for some reason
		// we've been asked to persist a historical / already-archived record.
		if !result.ArchivedAt.IsZero() {
			return tx.
				Create(result).
				Error
		}

		resultCurrent, err := s.getLatestResourceScan(tx, result.ResourceGroup, result.ResourceVersion, result.ResourceKind, result.ResourceNamespace, result.ResourceName)
		if err != nil {
			return err
		}

		// Calculate whether to persist this as a new result, archiving any old one, or whether to simply update
		// the checked_at time.
		newResult := resultCurrent == nil
		if resultCurrent != nil {
			newResult = resultCurrent.OverallStatus != result.OverallStatus || len(resultCurrent.Results) != len(result.Results)
			if !newResult {
				for i, rr := range result.Results {
					if resultCurrent.Results[i].Status != rr.Status || resultCurrent.Results[i].Message != rr.Message {
						newResult = true
						break
					}
				}
			}
		}

		if !newResult {
			// Not a new result, just update the checked_at on the existing result and we're done.
			return tx.Exec(
				"UPDATE security_scan_results SET checked_at = ? WHERE archived_at IS NULL AND resource_group = ? AND resource_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
				result.CheckedAt,
				result.ResourceGroup,
				result.ResourceVersion,
				result.ResourceKind,
				result.ResourceNamespace,
				result.ResourceName,
			).Error
		}

		if resultCurrent != nil {
			// We've got a new result and a current one, so archive the old one
			err := tx.Exec(
				"UPDATE security_scan_results SET archived_at = ? WHERE archived_at IS NULL AND resource_group = ? AND resource_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
				result.CheckedAt,
				result.ResourceGroup,
				result.ResourceVersion,
				result.ResourceKind,
				result.ResourceNamespace,
				result.ResourceName,
			).Error

			if err != nil {
				return err
			}
		}

		// Persist the new result.
		return tx.
			Create(result).
			Error
	})
}

func (s *securityImpl) GetOverview(ctx context.Context) (*model.SecurityOverview, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	return s.getOverview(ctx, "")
}

func (s *securityImpl) GetTeamOverview(ctx context.Context, team string) (*model.SecurityOverview, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	return s.getOverview(ctx, team)
}

func (s *securityImpl) getOverview(ctx context.Context, team string) (*model.SecurityOverview, error) {
	var result []struct {
		ResourceGroup     string
		ResourceVersion   string
		ResourceKind      string
		ResourceNamespace string
		ResourceName      string
		OverallStatus     string
		CheckedAt         time.Time
		RuleCode          string
		Status            string
	}

	query := s.conn.Table("security_scan_results").
		Joins("INNER JOIN security_rule_results ON security_scan_results.id = security_rule_results.scan_id").
		Select("security_scan_results.resource_group, security_scan_results.resource_version, security_scan_results.resource_kind, security_scan_results.resource_namespace, security_scan_results.resource_name, security_scan_results.checked_at, security_scan_results.overall_status, security_rule_results.rule_code, security_rule_results.status").
		Where("security_scan_results.archived_at IS NULL")

	if team != "" {
		query = query.
			Where("security_scan_results.owning_team = ?", team)
	}

	err := query.
		Order("security_scan_results.resource_group, security_scan_results.resource_version, security_scan_results.resource_kind, security_scan_results.resource_namespace, security_scan_results.resource_name").
		Scan(&result).
		Error

	if err != nil {
		return nil, err
	}

	overview := model.SecurityOverview{
		OpenIssueCounts: map[string]uint64{},
	}

	// Count totals and push results
	currRes := &model.SecurityResourceOverview{}
	for _, r := range result {
		if r.ResourceGroup != currRes.ResourceGroup || r.ResourceVersion != currRes.ResourceVersion || r.ResourceKind != currRes.ResourceKind || r.ResourceNamespace != currRes.ResourceNamespace || r.ResourceName != currRes.ResourceName {
			currRes = &model.SecurityResourceOverview{
				SecurityResourceReference: model.SecurityResourceReference{
					ResourceGroup:     r.ResourceGroup,
					ResourceVersion:   r.ResourceVersion,
					ResourceKind:      r.ResourceKind,
					ResourceNamespace: r.ResourceNamespace,
					ResourceName:      r.ResourceName,
				},
				LastChecked:     r.CheckedAt,
				OverallStatus:   r.OverallStatus,
				OpenIssueCounts: map[string]uint64{},
			}
			overview.Resources = append(overview.Resources, *currRes)
		}
		currRes.OpenIssueCounts[r.Status]++
		overview.OpenIssueCounts[r.Status]++
	}

	return &overview, nil
}
