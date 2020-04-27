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
	GetLatestResourceScan(ctx context.Context, apiVersion string, kind string, namespace string, name string) (*model.SecurityScanResult, error)
	// ListResourceScanHistory returns all scans for a specific resource
	ListResourceScanHistory(ctx context.Context, apiVersion string, kind string, namespace string, name string) ([]*model.SecurityScanResult, error)
	// StoreScan stores the result of a scan in the store. If the supplied result has a zero archived_at, any
	// previous scan with a zero archived_at for the same resource name/namespace to set archived_at to the
	// time of this scan
	StoreScan(context.Context, *model.SecurityScanResult) error
}

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
			"resource_api_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
			terms.GetAPIVersion(),
			terms.GetKind(),
			terms.GetNamespace(),
			terms.GetName(),
		)
	} else {
		if terms.HasTeam() {
			q = q.Where("owning_team = ?", terms.GetTeam())
		}
		if terms.HasAPIVersion() {
			q = q.Where("resource_api_version = ?", terms.GetAPIVersion())
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

func (s *securityImpl) GetLatestResourceScan(ctx context.Context, apiVersion string, kind string, namespace string, name string) (*model.SecurityScanResult, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	scan := &model.SecurityScanResult{}
	err := s.conn.Preload("Results").
		Where("resource_api_version = ?", apiVersion).
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

func (s *securityImpl) ListResourceScanHistory(ctx context.Context, apiVersion string, kind string, namespace string, name string) ([]*model.SecurityScanResult, error) {
	// This is just a convenience wrapper for the generic list scans function:
	return s.ListScans(ctx, false, Filter.WithIdentity(apiVersion, kind, namespace, name))
}

func (s *securityImpl) StoreScan(ctx context.Context, result *model.SecurityScanResult) error {
	// Ensure the archiving and inserting is done atomically in a single transaction.
	return s.conn.Transaction(func(tx *gorm.DB) error {

		// Archive any unarchived previous results for this name/namespace if ArchivedAt is
		// not set on this result (i.e. whenever we're not recording an already-archived
		// result)
		if result.ArchivedAt.IsZero() {
			err := tx.Exec(
				"UPDATE security_scan_results SET archived_at = ? WHERE archived_at IS NULL AND resource_api_version = ? AND resource_kind = ? AND resource_namespace = ? AND resource_name = ?",
				result.CheckedAt,
				result.ResourceAPIVersion,
				result.ResourceKind,
				result.ResourceNamespace,
				result.ResourceName,
			).Error
			if err != nil {
				return err
			}
		}

		// Now do the insert.
		return tx.
			Create(result).
			Error
	})
}
